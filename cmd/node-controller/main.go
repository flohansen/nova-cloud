package main

import (
	"fmt"
	"log"
	"net"

	"github.com/flohansen/nova-cloud/internal/node-controller/api"
	proto "github.com/flohansen/nova-cloud/proto/go"
	"google.golang.org/grpc"
)

func run() error {
	lis, err := net.Listen("tcp", ":3010")
	if err != nil {
		return fmt.Errorf("could not listen: %s", err)
	}

	srv := grpc.NewServer()
	proto.RegisterNodeControllerServer(srv, &api.Server{})

	return fmt.Errorf("could not serve: %s", srv.Serve(lis))
}

func main() {
	log.Fatalf("could not serve RPC server: %s", run())
}
