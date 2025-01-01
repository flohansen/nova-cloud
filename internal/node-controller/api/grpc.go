package api

import (
	"context"
	"fmt"
	"log"
	"net"

	proto "github.com/flohansen/nova-cloud/proto/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type PhyisicalNode struct {
	Addr string
}

func (n *PhyisicalNode) GetResources(ctx context.Context) (*proto.GetMachineInfoResponse, error) {
	conn, err := grpc.NewClient(n.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not create gRPC client: %s", err)
	}
	defer conn.Close()

	c := proto.NewCapacityServiceClient(conn)

	res, err := c.GetMachineInfo(ctx, &proto.GetMachineInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not get machine info: %s", err)
	}

	return res, nil
}

type Server struct {
	proto.UnimplementedNodeControllerServer
	nodes []PhyisicalNode
}

func (s *Server) RegisterNode(ctx context.Context, r *proto.RegisterNodeRequest) (*proto.Nothing, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get peer from context")
	}

	if p.Addr == nil {
		return nil, status.Errorf(codes.Internal, "peer has no address")
	}

	addr, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		return nil, status.Errorf(codes.Internal, "could not cast tcp address")
	}

	var addrString string
	if ipv4 := addr.IP.To4(); ipv4 != nil {
		addrString = fmt.Sprintf("%s:%d", addr.IP, r.Port)
	} else {
		addrString = fmt.Sprintf("[%s]:%d", addr.IP, r.Port)
	}

	log.Printf("new physical node: %s", addrString)
	s.nodes = append(s.nodes, PhyisicalNode{
		Addr: addrString,
	})

	return &proto.Nothing{}, nil
}

func (s *Server) ProvisionResources(ctx context.Context, r *proto.ProvisionResourcesRequest) (*proto.Nothing, error) {
	for _, node := range s.nodes {
		_, err := node.GetResources(ctx)
		if err != nil {
			log.Printf("could not fetch physical node info: %s", err)
		}

		// TODO: Add placement algorithm
	}

	return &proto.Nothing{}, nil
}
