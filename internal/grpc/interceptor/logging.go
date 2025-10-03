package interceptor

import (
	"context"

	"github.com/flohansen/nova-cloud/internal/logging"
	"google.golang.org/grpc"
)

func UnaryLogger(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = logging.WithContext(ctx, logger)
		return handler(ctx, req)
	}
}

func UnaryRequestLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log := logging.FromContext(ctx)
		log.Info("incoming request", "rpc", info.FullMethod)

		res, err := handler(ctx, req)
		if err != nil {
			log.Error("request error", "error", err)
		}

		log.Info("responding", "rpc", info.FullMethod)
		return res, err
	}
}

func StreamLogger(logger logging.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ss = &wrappedStream{
			ServerStream: ss,
			ctx:          logging.WithContext(ss.Context(), logger),
		}
		return handler(srv, ss)
	}
}

func StreamRequestLogging() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := logging.FromContext(ss.Context())
		log.Info("incoming request", "rpc", info.FullMethod)

		err := handler(srv, ss)
		if err != nil {
			log.Error("request error", "error", err)
		}

		log.Info("responding", "rpc", info.FullMethod)
		return err
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ws *wrappedStream) Context() context.Context {
	return ws.ctx
}
