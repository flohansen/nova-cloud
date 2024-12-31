package api

import (
	"context"
	"log"
	"net"

	proto "github.com/flohansen/nova-cloud/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedNodeControllerServer
}

func (s *Server) RegisterNode(ctx context.Context, r *proto.RegisterNodeRequest) (*proto.Nothing, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get peer from context")
	}

	if p.Addr != nil {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP := addr.IP
			log.Printf("new register request from %s", clientIP)
		}
	}

	return &proto.Nothing{}, nil
}
