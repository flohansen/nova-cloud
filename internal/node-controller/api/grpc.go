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

	c := proto.NewNodeClient(conn)

	res, err := c.GetMachineInfo(ctx, &proto.GetMachineInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not get machine info: %s", err)
	}

	return res, nil
}

func (n *PhyisicalNode) Aquire(ctx context.Context, r *proto.AquireRequest) error {
	conn, err := grpc.NewClient(n.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not create gRPC client: %s", err)
	}
	defer conn.Close()

	c := proto.NewNodeClient(conn)

	_, err = c.Aquire(ctx, r)
	if err != nil {
		return fmt.Errorf("could not aquire resources: %s", err)
	}

	return nil
}

type Server struct {
	proto.UnimplementedNodeControllerServer
	nodes []PhyisicalNode
}

func (s *Server) RegisterNode(ctx context.Context, r *proto.RegisterNodeRequest) (*proto.RegisterNodeResponse, error) {
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

	return &proto.RegisterNodeResponse{}, nil
}

func (s *Server) ProvisionResources(ctx context.Context, r *proto.ProvisionResourcesRequest) (*proto.ProvisionResourcesResponse, error) {
	for _, node := range s.nodes {
		res, err := node.GetResources(ctx)
		if err != nil {
			log.Printf("could not fetch physical node info: %s", err)
			continue
		}

		if res.CpuCores < r.CpuCores {
			continue
		}

		if res.MemoryBytes < r.MemoryBytes {
			continue
		}

		if err := node.Aquire(ctx, &proto.AquireRequest{
			CpuCores:    r.CpuCores,
			MemoryBytes: r.MemoryBytes,
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "invalid node configuration: %s", err)
		}

		return &proto.ProvisionResourcesResponse{}, nil
	}

	return nil, status.Errorf(codes.Internal, "no node has enough resources for this request %+v", r)
}
