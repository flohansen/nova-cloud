package handler

import (
	"context"
	"runtime"

	v1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"google.golang.org/grpc"
)

var _ v1.NodeAgentServiceServer = &NodeAgentHandler{}

type NodeAgentHandler struct {
	v1.UnsafeNodeAgentServiceServer
}

func NewNodeAgentHandler() *NodeAgentHandler {
	return &NodeAgentHandler{}
}

func (n *NodeAgentHandler) Desc() *grpc.ServiceDesc {
	return &v1.NodeAgentService_ServiceDesc
}

func (n *NodeAgentHandler) GetResources(ctx context.Context, req *v1.GetResourcesRequest) (*v1.GetResourcesResponse, error) {
	cpuCores := int32(runtime.NumCPU())

	cpuArch := v1.CpuArch_CPU_ARCH_UNSPECIFIED
	switch runtime.GOARCH {
	case "386":
		cpuArch = v1.CpuArch_CPU_ARCH_X86
	case "amd64":
		cpuArch = v1.CpuArch_CPU_ARCH_X86_64
	case "arm":
		cpuArch = v1.CpuArch_CPU_ARCH_ARM
	case "arm64":
		cpuArch = v1.CpuArch_CPU_ARCH_ARM64
	}

	return &v1.GetResourcesResponse{
		CpuCores: cpuCores,
		CpuArch:  cpuArch,
	}, nil
}
