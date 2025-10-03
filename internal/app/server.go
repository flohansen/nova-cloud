package app

import (
	"context"
	"net"

	"github.com/flohansen/nova-cloud/internal/logging"
	v1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:generate mockgen -destination=mocks/node_service_server.go -package=mocks github.com/flohansen/nova-cloud/internal/proto/novacloud/v1 NodeAgentServiceServer
//go:generate mockgen -destination=mocks/logger.go -package=mocks github.com/flohansen/nova-cloud/internal/logging Logger

type Server struct {
	logger           logging.Logger
	listener         net.Listener
	controller       v1.NodeAgentServiceServer
	server           *grpc.Server
	enableReflection bool
}

func NewServer(listener net.Listener, controller v1.NodeAgentServiceServer, logger logging.Logger, opts ...ServerOpt) *Server {
	s := &Server{
		logger:           logger,
		listener:         listener,
		controller:       controller,
		enableReflection: false,
	}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	s.server = grpc.NewServer()

	if s.enableReflection {
		s.logger.Info("enabling gRPC reflection")
		reflection.Register(s.server)
	}

	v1.RegisterNodeAgentServiceServer(s.server, s.controller)

	go s.handleShutdown(ctx)

	s.logger.Info("starting gRPC server", "addr", s.listener.Addr().String())
	return s.server.Serve(s.listener)
}

func (s *Server) handleShutdown(ctx context.Context) {
	<-ctx.Done()
	s.logger.Info("context done, shutting down server")
	s.server.GracefulStop()
}

type ServerOpt func(*Server)

func WithReflection() ServerOpt {
	return func(s *Server) {
		s.enableReflection = true
	}
}
