package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/flohansen/nova-cloud/internal/node/api"
	proto "github.com/flohansen/nova-cloud/proto/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port           = flag.Int("port", 3000, "The port on which the node server should listen")
	controllerAddr = flag.String("controller-addr", "localhost:3010", "The address of the node controller")
)

func main() {
	flag.Parse()

	if err := registerNode(); err != nil {
		log.Fatalf("could not register node: %s", err)
	}

	if err := serveMetrics(); err != nil {
		log.Fatalf("could not serve metrics: %s", err)
	}
}

func registerNode() error {
	conn, err := grpc.NewClient(*controllerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could create gRPC client: %s", err)
	}
	defer conn.Close()

	c := proto.NewNodeControllerClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = c.RegisterNode(ctx, &proto.RegisterNodeRequest{
		Port: int32(*port),
	})
	if err != nil {
		return fmt.Errorf("could not register node: %s", err)
	}

	return nil
}

func serveMetrics() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("could not listen: %s", err)
	}

	srv := grpc.NewServer()
	proto.RegisterNodeServer(srv, &api.Server{})

	return fmt.Errorf("could not serve: %s", srv.Serve(lis))
}
