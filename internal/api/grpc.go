package api

import (
	"context"

	"github.com/flohansen/nova-cloud/capacity/internal/proc"
	proto "github.com/flohansen/nova-cloud/proto/go"
)

type Server struct {
	proto.UnimplementedCapacityServiceServer
}

func (s *Server) GetMachineInfo(context.Context, *proto.GetMachineInfoRequest) (*proto.GetMachineInfoResponse, error) {
	cores := proc.GetCores()
	memInfo := proc.GetMemInfo()

	return &proto.GetMachineInfoResponse{
		Cpu: cores,
		Ram: memInfo.MemTotal,
	}, nil
}
