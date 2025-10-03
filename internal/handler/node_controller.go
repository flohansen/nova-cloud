package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	novacloudv1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"github.com/flohansen/nova-cloud/sql/generated/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type NodeControllerHandler struct {
	novacloudv1.UnimplementedNodeControllerServiceServer

	q *database.Queries
}

func NewNodeControllerHandler(db database.DBTX) *NodeControllerHandler {
	h := &NodeControllerHandler{
		q: database.New(db),
	}

	go h.startHealthChecks(context.Background())
	return h
}

func (h *NodeControllerHandler) RegisterNode(ctx context.Context, req *novacloudv1.RegisterNodeRequest) (*novacloudv1.RegisterNodeResponse, error) {
	p, _ := peer.FromContext(ctx)

	tcpAddr, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		return nil, status.Error(codes.Internal, "peer address is not a TCP address")
	}

	target := fmt.Sprintf("%s:%d", tcpAddr.IP, req.Port)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create grpc client: %v", err)
	}

	client := novacloudv1.NewNodeAgentServiceClient(conn)
	res, err := client.GetResources(ctx, &novacloudv1.GetResourcesRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "node get resources error: %v", err)
	}

	if err := h.q.UpsertNode(ctx, database.UpsertNodeParams{
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
	nodes, err := h.q.GetNodes(stream.Context())
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

func (h *NodeControllerHandler) startHealthChecks(ctx context.Context) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	for {
		time.Sleep(time.Second)

		nodes, err := h.q.GetNodes(ctx)
		if err != nil {
			log.Error("get nodes", "error", err)
			continue
		}

		for _, node := range nodes {
			target := fmt.Sprintf("%s:%d", node.Ip, node.Port)
			conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error("new grpc client", "error", err)
				continue
			}

			client := novacloudv1.NewNodeAgentServiceClient(conn)
			_, err = client.GetResources(ctx, &novacloudv1.GetResourcesRequest{})
			if err != nil {
				log.Info("healthcheck failed", "error", err, "target", target)

				if err := h.q.DeleteNode(ctx, target); err != nil {
					log.Error("delete node", "error", err, "target", target)
				}

				continue
			}
		}
	}
}
