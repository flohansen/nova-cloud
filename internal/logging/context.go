package logging

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const (
	contextKeyLogger contextKey = "logger"
)

func WithContext(parent context.Context, logger Logger) context.Context {
	return context.WithValue(parent, contextKeyLogger, logger)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(contextKeyLogger).(Logger)
	if !ok {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	return logger
}
