package doubles

import (
	"context"
	"net"
	"testing"

	novacloudv1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"google.golang.org/grpc"
)

type TestNodeAgent struct {
	novacloudv1.UnimplementedNodeAgentServiceServer

	Addr     net.Addr
	CpuCores int32
	CpuArch  novacloudv1.CpuArch
}

func StartTestNodeAgent(t *testing.T, cpus int32, arch novacloudv1.CpuArch) *TestNodeAgent {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(lis.Addr().String())

	srv := grpc.NewServer()
	novacloudv1.RegisterNodeAgentServiceServer(srv, &TestNodeAgent{
		CpuCores: cpus,
		CpuArch:  arch,
	})
	go func() {
		srv.Serve(lis)
	}()
	t.Cleanup(func() {
		srv.GracefulStop()
	})

	return &TestNodeAgent{
		Addr:     lis.Addr(),
		CpuCores: cpus,
		CpuArch:  arch,
	}
}

func (t *TestNodeAgent) GetResources(context.Context, *novacloudv1.GetResourcesRequest) (*novacloudv1.GetResourcesResponse, error) {
	return &novacloudv1.GetResourcesResponse{
		CpuCores: t.CpuCores,
		CpuArch:  t.CpuArch,
	}, nil
}
