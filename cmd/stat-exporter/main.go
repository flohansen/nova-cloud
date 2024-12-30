package main

import (
	"log"
	"net"

	"github.com/flohansen/nova-cloud/capacity/internal/api"
	proto "github.com/flohansen/nova-cloud/proto/go"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("could not listen: %s", err)
	}

	srv := grpc.NewServer()
	proto.RegisterCapacityServiceServer(srv, &api.Server{})
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("could not serve: %s", err)
	}
}
