package app_test

import (
	"context"
	"os/signal"
	"syscall"
	"testing"

	"github.com/flohansen/nova-cloud/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestSignalContext(t *testing.T) {
	t.Run("should cancel context on signals", func(t *testing.T) {
		// assign
		tests := []struct {
			name   string
			signal syscall.Signal
		}{
			{name: "SIGINT", signal: syscall.SIGINT},
			{name: "SIGTERM", signal: syscall.SIGTERM},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				signal.Reset()

				// act
				ctx := app.SignalContext()
				syscall.Kill(syscall.Getpid(), tt.signal)

				// assert
				<-ctx.Done()
				assert.Equal(t, context.Canceled, ctx.Err())
			})
		}
	})
}
