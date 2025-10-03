package handler

import (
	"context"
	"fmt"
	"net"

	"github.com/flohansen/nova-cloud/internal/domain"
	novacloudv1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type NodeRepository interface {
	FindAll(ctx context.Context) ([]domain.Node, error)
	CreateOrUpdate(ctx context.Context, node domain.Node) error
	Delete(ctx context.Context, nodeID string) error
}

//go:generate mockgen -destination=mocks/node_stream.go -package=mocks github.com/flohansen/nova-cloud/internal/proto/novacloud/v1 NodeControllerService_GetNodesServer

type NodeControllerHandler struct {
	novacloudv1.UnimplementedNodeControllerServiceServer

	nodeRepo NodeRepository
}

func NewNodeControllerHandler(nodeRepo NodeRepository) *NodeControllerHandler {
	return &NodeControllerHandler{
		nodeRepo: nodeRepo,
	}
}

func (h *NodeControllerHandler) RegisterNode(ctx context.Context, req *novacloudv1.RegisterNodeRequest) (*novacloudv1.RegisterNodeResponse, error) {
	p, _ := peer.FromContext(ctx)

	tcpAddr, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		return nil, status.Error(codes.Internal, "peer address is not a TCP address")
	}

	var target string
	if ip := tcpAddr.IP.To4(); ip != nil {
		target = fmt.Sprintf("%s:%d", tcpAddr.IP, req.Port)
	} else {
		target = fmt.Sprintf("[%s]:%d", tcpAddr.IP, req.Port)
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create grpc client: %v", err)
	}

	client := novacloudv1.NewNodeAgentServiceClient(conn)
	res, err := client.GetResources(ctx, &novacloudv1.GetResourcesRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "node get resources error: %v", err)
	}

	if err := h.nodeRepo.CreateOrUpdate(ctx, domain.Node{
		Ip:      tcpAddr.IP.String(),
		NodeID:  target,
		Port:    int64(req.Port),
		Cpus:    int64(res.CpuCores),
		CpuArch: int64(res.CpuArch),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "upsert node info error: %v", err)
	}

	return &novacloudv1.RegisterNodeResponse{}, nil
}

func (h *NodeControllerHandler) GetNodes(req *novacloudv1.GetNodesRequest, stream grpc.ServerStreamingServer[novacloudv1.GetNodesResponse]) error {
	nodes, err := h.nodeRepo.FindAll(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "get nodes error: %v", err)
	}

	for _, node := range nodes {
		stream.Send(&novacloudv1.GetNodesResponse{
			Ip:       node.Ip,
			Port:     int32(node.Port),
			CpuArch:  novacloudv1.CpuArch(node.CpuArch),
			CpuCores: int32(node.Cpus),
		})
	}

	return nil
}
