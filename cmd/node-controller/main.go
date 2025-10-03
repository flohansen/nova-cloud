package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"

	"github.com/flohansen/nova-cloud/internal/app"
	"github.com/flohansen/nova-cloud/internal/grpc"
	"github.com/flohansen/nova-cloud/internal/handler"
	"github.com/flohansen/nova-cloud/internal/logging"
	"github.com/flohansen/nova-cloud/internal/repository"
	"github.com/flohansen/nova-cloud/sql/migrations"
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

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Error("sql open error", "error", err)
		os.Exit(1)
	}

	if err := migrations.Run(db, "novacloud"); err != nil {
		log.Error("migration error", "error", err)
		os.Exit(1)
	}

	nodeRepo := repository.NewNodeRepository(db)
	nodeControllerService := handler.NewNodeControllerHandler(nodeRepo)

	if err := grpc.NewServer(
		grpc.WithListenAddr(config.ListenAddr),
		grpc.WithService(nodeControllerService),
	).Serve(ctx); err != nil {
		log.Error("serve error", "error", err)
		os.Exit(1)
	}
}
