package grpc

import (
	"context"

	"github.com/enesanbar/go-service/core/errors"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ServerOptionUnaryInterceptorErrorHandlerParams holds the parameters for creating a gRPC server option for unary interceptor error handling.
type ServerOptionUnaryInterceptorErrorHandlerParams struct {
	fx.In

	Logger log.Factory
}

// NewServerOptionUnaryInterceptorErrorHandler creates a new gRPC server option for unary interceptor error handling.
func NewServerOptionUnaryInterceptorErrorHandler(p ServerOptionUnaryInterceptorErrorHandlerParams) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		m, err := handler(ctx, req)
		if err != nil {
			// st, _ := status.FromError(err)
			// return m, st.Err()
			code := ErrorStatus(err)
			message := errors.ErrorMessage(err)
			// data := customErr.ErrorData(err)
			p.Logger.For(ctx).With(zap.Error(err)).Error("gRPC unary interceptor error")

			return m, status.Error(code, message)
		}

		return m, err
	})
}
