package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/flohansen/nova-cloud/internal/app"
	"github.com/flohansen/nova-cloud/internal/grpc"
	"github.com/flohansen/nova-cloud/internal/handler"
	"github.com/flohansen/nova-cloud/internal/logging"
)

type config struct {
	ListenAddr string
}

func main() {
	var config config
	flag.StringVar(&config.ListenAddr, "listen", "0.0.0.0:5050", "The listen address of the gRPC server")
	flag.Parse()

	ctx := app.SignalContext()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	ctx = logging.WithContext(ctx, log)

	nodeAgentService := handler.NewNodeAgentHandler()

	if err := grpc.NewServer(
		grpc.WithListenAddr(config.ListenAddr),
		grpc.WithService(nodeAgentService),
	).Serve(ctx); err != nil {
		log.Error("serve error", "error", err)
		os.Exit(1)
	}
}
