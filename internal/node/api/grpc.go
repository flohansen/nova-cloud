package api

import (
	"context"
	"log"

	"github.com/flohansen/nova-cloud/internal/node/proc"
	proto "github.com/flohansen/nova-cloud/proto/go"
)

type Server struct {
	proto.UnimplementedNodeServer
}

func (s *Server) GetMachineInfo(context.Context, *proto.GetMachineInfoRequest) (*proto.GetMachineInfoResponse, error) {
	cores := proc.GetCores()
	memInfo := proc.GetMemInfo()

	return &proto.GetMachineInfoResponse{
		CpuCores:    cores,
		MemoryBytes: memInfo.MemTotal * 1024,
	}, nil
}

func (s *Server) Aquire(ctx context.Context, r *proto.AquireRequest) (*proto.AquireResponse, error) {
	log.Printf("starting new VM with %+v", r)
	return &proto.AquireResponse{}, nil
}
