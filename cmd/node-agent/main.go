package main

import (
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/flohansen/nova-cloud/internal/app"
	"github.com/flohansen/nova-cloud/internal/controller"
)

type flags struct {
	EnableReflection bool
	ListenAddr       string
}

func main() {
	var flags flags
	flag.BoolVar(&flags.EnableReflection, "reflection", false, "If the gRPC server should enable reflection")
	flag.StringVar(&flags.ListenAddr, "listen", "0.0.0.0:5050", "The listen address of the gRPC server")
	flag.Parse()

	ctx := app.SignalContext()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	lis, err := net.Listen("tcp", flags.ListenAddr)
	if err != nil {
		logger.Error("could not create network listener", "error", err)
		os.Exit(1)
	}
	defer lis.Close()

	var opts []app.ServerOpt
	if flags.EnableReflection {
		opts = append(opts, app.WithReflection())
	}

	controller := controller.NewNodeController()
	srv := app.NewServer(lis, controller, logger, opts...)
	if err := srv.Run(ctx); err != nil {
		logger.Error("gRPC server error", "error", err)
		os.Exit(1)
	}
}
