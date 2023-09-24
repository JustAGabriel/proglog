package server

import (
	"context"

	api "github.com/justagabriel/proglog/api/v1"
	"google.golang.org/grpc"
)

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}

type Config struct {
	CommitLog CommitLog
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

func (s *grpcServer) Create(ctx context.Context, req *api.CreateRecordRequest) (*api.CreateRecordResponse, error) {
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
	gsrv := grpc.NewServer(opts...)
	srv, err := newGRPCServer(config)
	if err != nil {
		return nil, err
	}

	api.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}
