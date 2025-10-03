package app_test

import (
	"context"
	"testing"

	"github.com/flohansen/nova-cloud/internal/app"
	"github.com/flohansen/nova-cloud/internal/app/mocks"
	v1 "github.com/flohansen/nova-cloud/internal/proto/novacloud/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/test/bufconn"
)

func TestServer_NewServer(t *testing.T) {
	// assign
	// act
	srv := app.NewServer(nil, nil, nil)

	// assert
	assert.NotNil(t, srv)
}

func TestServer_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerMock := &controllerMockWrapper{
		mock: mocks.NewMockNodeAgentServiceServer(ctrl),
	}
	loggerMock := mocks.NewMockLogger(ctrl)

	t.Run("should write logs about starting and terminating server", func(t *testing.T) {
		// assign
		ctx, cancel := context.WithCancel(context.Background())
		listenerStub := bufconn.Listen(1024 * 1024)
		defer listenerStub.Close()

		srv := app.NewServer(listenerStub, controllerMock, loggerMock)

		loggerMock.EXPECT().
			Info("starting gRPC server", "addr", listenerStub.Addr().String()).
			Do(func(msg string, v ...any) { cancel() }).
			Times(1)
		loggerMock.EXPECT().
			Info("context done, shutting down server").
			Times(1)

		// act
		err := make(chan error)
		go func() {
			defer close(err)
			err <- srv.Run(ctx)
		}()

		// assert
		assert.NoError(t, <-err)
	})

	t.Run("should write log about reflection is enabled", func(t *testing.T) {
		// assign
		ctx, cancel := context.WithCancel(context.Background())
		listenerStub := bufconn.Listen(1024 * 1024)
		defer listenerStub.Close()

		srv := app.NewServer(listenerStub, controllerMock, loggerMock, app.WithReflection())

		loggerMock.EXPECT().
			Info("enabling gRPC reflection").
			Times(1)
		loggerMock.EXPECT().
			Info("starting gRPC server", "addr", listenerStub.Addr().String()).
			Do(func(msg string, v ...any) { cancel() }).
			Times(1)
		loggerMock.EXPECT().
			Info("context done, shutting down server").
			Times(1)

		// act
		err := make(chan error)
		go func() {
			defer close(err)
			err <- srv.Run(ctx)
		}()

		// assert
		assert.NoError(t, <-err)
	})
}

type controllerMockWrapper struct {
	v1.UnimplementedNodeAgentServiceServer
	mock *mocks.MockNodeAgentServiceServer
}

func (m *controllerMockWrapper) GetResources(ctx context.Context, req *v1.GetResourcesRequest) (*v1.GetResourcesResponse, error) {
	return m.mock.GetResources(ctx, req)
}
