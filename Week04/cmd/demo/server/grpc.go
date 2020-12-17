package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server

	addr string
	opts []grpc.ServerOption
}

func NewServer(addr string, opts ...grpc.ServerOption) *Server {
	srv := &Server{
		addr: addr,
		opts: opts,
	}
	srv.Server = grpc.NewServer(opts...)
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	return nil
}
