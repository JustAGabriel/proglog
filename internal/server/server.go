package server

import (
	"context"

	api "github.com/justagabriel/proglog/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	objectWildcard string = "*"
	createAction   string = "create"
	getAction      string = "get"
)

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}

type Authorizer interface {
	Authorize(subject, object, action string) error
}

type Config struct {
	CommitLog  CommitLog
	Authorizer Authorizer
}

type grpcServer struct {
	api.UnimplementedLogServer
	*Config
}

func newGRPCServer(config *Config) (*grpcServer, error) {
	srv := &grpcServer{
		Config: config,
	}
	return srv, nil
}

var _ api.LogServer = (*grpcServer)(nil)

type subjectContextKey struct{}

func authenticate(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.New(codes.Unknown, "couldn't find peer info").Err()
	}

	if peer.AuthInfo == nil {
		return context.WithValue(ctx, subjectContextKey{}, ""), nil
	}

	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
	subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	ctx = context.WithValue(ctx, subjectContextKey{}, subject)
	return ctx, nil
}

func subject(ctx context.Context) string {
	return ctx.Value(subjectContextKey{}).(string)
}

func (s *grpcServer) Create(ctx context.Context, req *api.CreateRecordRequest) (*api.CreateRecordResponse, error) {
	if err := s.Authorizer.Authorize(subject(ctx), objectWildcard, createAction); err != nil {
		return nil, err
	}
	offset, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &api.CreateRecordResponse{Offset: offset}, nil
}

func (s *grpcServer) Get(ctx context.Context, req *api.GetRecordRequest) (*api.GetRecordResponse, error) {
	rec, err := s.CommitLog.Read(req.GetOffset())
	if err != nil {
		return nil, err
	}

	return &api.GetRecordResponse{Record: rec}, nil
}

func (s *grpcServer) CreateStream(stream api.Log_CreateStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		res, err := s.Create(stream.Context(), req)
		if err != nil {
			return err
		}

		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

func (s *grpcServer) GetStream(stream api.Log_GetStreamServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			req, err := stream.Recv()
			if err != nil {
				return err
			}

			res, err := s.Get(stream.Context(), req)
			switch err.(type) {
			case nil:
			case api.ErrOffsetOutOfRange:
				continue
			default:
				return err
			}

			err = stream.Send(res)
			if err != nil {
				return err
			}
			req.Offset++
		}
	}
}

func NewGRPCServer(config *Config, opts ...grpc.ServerOption) (*grpc.Server, error) {
	// todo: implement auth middleware usage!
	authMiddleware := grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(grpc_auth.StreamServerINterceptor(authenticate())),
	)
	gsrv := grpc.NewServer(opts...)
	srv, err := newGRPCServer(config)
	if err != nil {
		return nil, err
	}

	api.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}
