package handler_test

import (
	"net"
	"testing"

	"github.com/flohansen/nova-cloud/internal/domain"
	"github.com/flohansen/nova-cloud/internal/handler"
	"github.com/flohansen/nova-cloud/internal/handler/mocks"
	novacloudv1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"github.com/flohansen/nova-cloud/testing/doubles"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/peer"
)

func TestNodeControllerHandler_RegisterNode(t *testing.T) {
	// given
	agent := doubles.StartTestNodeAgent(t, 4, novacloudv1.CpuArch_CPU_ARCH_ARM)
	nodeRepo := doubles.NewTestNodeRepository()
	h := handler.NewNodeControllerHandler(nodeRepo)

	// when
	ctx := peer.NewContext(t.Context(), &peer.Peer{
		Addr: agent.Addr,
	})
	res, err := h.RegisterNode(ctx, &novacloudv1.RegisterNodeRequest{
		Port: int32(agent.Addr.(*net.TCPAddr).Port),
	})

	// then
	assert.NoError(t, err)
	assert.NotNil(t, res)
	node, ok := nodeRepo.Nodes[agent.Addr.String()]
	assert.True(t, ok, "node should be registered")
	assert.Equal(t, int64(4), node.Cpus)
	assert.Equal(t, int64(novacloudv1.CpuArch_CPU_ARCH_ARM), node.CpuArch)
	assert.Equal(t, agent.Addr.(*net.TCPAddr).IP.String(), node.Ip)
	assert.Equal(t, int64(agent.Addr.(*net.TCPAddr).Port), node.Port)
}

func TestNodeControllerHandler_GetNodes(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	nodes := []domain.Node{
		{NodeID: "node1"},
		{NodeID: "node2"},
		{NodeID: "node3"},
	}
	nodeRepo := doubles.NewTestNodeRepository().WithNodes(nodes...)
	h := handler.NewNodeControllerHandler(nodeRepo)
	req := &novacloudv1.GetNodesRequest{}
	stream := mocks.NewMockNodeControllerService_GetNodesServer[novacloudv1.GetNodesResponse](ctrl)

	stream.EXPECT().Context().Return(t.Context())

	for _, node := range nodes {
		stream.EXPECT().Send(&novacloudv1.GetNodesResponse{
			Ip:       node.Ip,
			Port:     int32(node.Port),
			CpuCores: int32(node.Cpus),
			CpuArch:  novacloudv1.CpuArch(node.CpuArch),
		}).Times(1)
	}

	// when
	err := h.GetNodes(req, stream)

	// then
	assert.NoError(t, err)
}

func TestNodeControllerHandler_CreateInstance(t *testing.T) {
	for _, tt := range []struct {
		name string
	}{
		{name: "no nodes"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := t.Context()
			req := &novacloudv1.CreateInstanceRequest{
				Vcpu: 1,
				Arch: novacloudv1.CpuArch_CPU_ARCH_ARM,
			}

			nodeRepo := doubles.NewTestNodeRepository()
			h := handler.NewNodeControllerHandler(nodeRepo)

			// when
			res, err := h.CreateInstance(ctx, req)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}
