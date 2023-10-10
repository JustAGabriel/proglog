package server

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	api "github.com/justagabriel/proglog/api/v1"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	createAction string = "create"
	getAction    string = "get"
)

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}

type Authorizer interface {
	Authorize(subject, action string) error
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
	subject := subject(ctx)
	err := s.Authorizer.Authorize(subject, getAction)
	if err != nil {
		return nil, err
	}
	offset, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &api.CreateRecordResponse{Offset: offset}, nil
}

func (s *grpcServer) Get(ctx context.Context, req *api.GetRecordRequest) (*api.GetRecordResponse, error) {
	subject := subject(ctx)
	err := s.Authorizer.Authorize(subject, getAction)
	if err != nil {
		return nil, err
	}
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

	logger := zap.L().Named("server")
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(
			func(duration time.Duration) zapcore.Field {
				return zap.Int64(
					"grpc.time_ns",
					duration.Nanoseconds(),
				)
			},
		),
	}

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	err := view.Register(ocgrpc.DefaultServerViews...)
	if err != nil {
		return nil, err
	}

	grpcOpts := []grpc.ServerOption{
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_zap.StreamServerInterceptor(logger, zapOpts...),
				grpc_auth.StreamServerInterceptor(authenticate),
			),
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
			grpc_auth.UnaryServerInterceptor(authenticate),
		)),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}

	opts = append(opts, grpcOpts...)

	gsrv := grpc.NewServer(opts...)
	srv, err := newGRPCServer(config)
	if err != nil {
		return nil, err
	}

	api.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}
