package controller

import (
	"context"
	"runtime"

	v1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
)

var _ v1.NodeServiceServer = &NodeController{}

type NodeController struct {
	v1.UnsafeNodeServiceServer
}

func NewNodeController() *NodeController {
	return &NodeController{}
}

func (n *NodeController) GetResources(context.Context, *v1.GetResourcesRequest) (*v1.GetResourcesResponse, error) {
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
		CpuCores:        cpuCores,
		CpuArchitecture: cpuArch,
	}, nil
}
