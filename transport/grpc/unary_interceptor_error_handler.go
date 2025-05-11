package grpc

import (
	"context"

	"github.com/enesanbar/go-service/errors"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GRPCServerOptionUnaryInterceptorErrorHandlerParams struct {
	fx.In

	Logger log.Factory
}

func NewGRPCServerOptionUnaryInterceptorErrorHandler(p GRPCServerOptionUnaryInterceptorErrorHandlerParams) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		m, err := handler(ctx, req)
		if err != nil {
			// st, _ := status.FromError(err)
			// return m, st.Err()
			code := errors.ErrorStatusGRPC(err)
			message := errors.ErrorMessage(err)
			// data := customErr.ErrorData(err)
			p.Logger.For(ctx).With(zap.Error(err)).Error("gRPC unary interceptor error")

			return m, status.Error(code, message)
		}

		return m, err
	})
}
