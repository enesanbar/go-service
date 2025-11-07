package grpc

import (
	"context"
	"errors"

	serviceErr "github.com/enesanbar/go-service/core/errors"
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

// NewUnaryServerInterceptorErrorHandler creates a new gRPC server option for unary interceptor error handling.
func NewUnaryServerInterceptorErrorHandler(p ServerOptionUnaryInterceptorErrorHandlerParams) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		m, err := handler(ctx, req)
		if err == nil {
			return m, nil
		}

		// If the error is already a gRPC status error, return it as is
		if statusErr, ok := status.FromError(err); ok {
			// log the error with details
			details := statusErr.Details()
			p.Logger.For(ctx).With(
				zap.String("method", info.FullMethod),
				zap.String("code", statusErr.Code().String()),
				zap.String("message", statusErr.Message()),
				zap.Any("details", details),
			).Error("gRPC unary interceptor status error")
			return m, err
		}

		code := ErrorStatus(err)
		message := serviceErr.ErrorMessage(err)

		//  TODO: Print only critical errors in ERROR level, others in WARN level
		p.Logger.For(ctx).With(
			zap.Error(errors.Unwrap(err)),
		).Error("gRPC unary interceptor error", zap.Error(err))

		return m, status.Error(code, message)
	}
}
