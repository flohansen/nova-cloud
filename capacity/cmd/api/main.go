package main

import (
	"context"
	"log"
	"net"

	"github.com/flohansen/nova-cloud/capacity/internal/proc"
	proto "github.com/flohansen/nova-cloud/proto/go"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedCapacityServiceServer
}

func (s *server) GetMachineInfo(context.Context, *proto.GetMachineInfoRequest) (*proto.GetMachineInfoResponse, error) {
	cores := proc.GetCores()
	memInfo := proc.GetMemInfo()

	return &proto.GetMachineInfoResponse{
		Cpu: cores,
		Ram: memInfo.MemTotal,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("could not listen: %s", err)
	}

	srv := grpc.NewServer()
	proto.RegisterCapacityServiceServer(srv, &server{})
	srv.Serve(lis)
}
