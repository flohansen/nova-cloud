package main

import (
	"database/sql"
	"log/slog"
	"net"
	"os"

	_ "modernc.org/sqlite"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/flohansen/nova-cloud/internal/handler"
	novacloudv1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"github.com/flohansen/nova-cloud/sql/migrations"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Error("sql open error", "error", err)
		os.Exit(1)
	}

	if err := migrations.Run(db, "novacloud"); err != nil {
		log.Error("migration error", "error", err)
		os.Exit(1)
	}

	nodeController := handler.NewNodeControllerHandler(db)

	srv := grpc.NewServer()
	novacloudv1.RegisterNodeControllerServiceServer(srv, nodeController)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Error("listen error", "error", err)
		os.Exit(1)
	}

	log.Info("server started", "addr", ":5050")
	if err := srv.Serve(lis); err != nil {
		log.Error("serve error", "error", err)
		os.Exit(1)
	}
}
