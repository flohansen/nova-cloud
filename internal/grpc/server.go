package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/flohansen/nova-cloud/internal/grpc/interceptor"
	"github.com/flohansen/nova-cloud/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service interface {
	Desc() *grpc.ServiceDesc
}

type Server struct {
	addr     string
	services []Service
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		addr: ":5050",
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Serve(ctx context.Context) error {
	log := logging.FromContext(ctx)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryLogger(log),
			interceptor.UnaryRequestLogging(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamLogger(log),
			interceptor.StreamRequestLogging(),
		),
	)
	for _, service := range s.services {
		srv.RegisterService(service.Desc(), service)
	}
	reflection.Register(srv)

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	srvErrs := make(chan error, 1)
	go func() {
		if err := srv.Serve(lis); err != nil {
			srvErrs <- err
		}
	}()

	log.Info("server started", "addr", s.addr)

	select {
	case <-ctx.Done():
		close(srvErrs)
		srv.GracefulStop()
		log.Info("server shutdown completed")
		return nil
	case err := <-srvErrs:
		return fmt.Errorf("server error: %w", err)
	}
}

type ServerOption func(*Server)

func WithListenAddr(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithService(service Service) ServerOption {
	return func(s *Server) {
		s.services = append(s.services, service)
	}
}
