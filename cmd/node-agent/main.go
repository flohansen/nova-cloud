package main

import (
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/flohansen/nova-cloud/internal/handler"
	novacloudv1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type config struct {
	ListenAddr string
}

func main() {
	var config config
	flag.StringVar(&config.ListenAddr, "listen", "0.0.0.0:5050", "The listen address of the gRPC server")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	controller := handler.NewNodeAgentHandler()

	srv := grpc.NewServer()
	novacloudv1.RegisterNodeAgentServiceServer(srv, controller)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		log.Error("could not create network listener", "error", err)
		os.Exit(1)
	}

	log.Info("server started", "addr", config.ListenAddr)
	if err := srv.Serve(lis); err != nil {
		log.Error("serve error", "error", err)
		os.Exit(1)
	}
}
