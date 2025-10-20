package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func panicHandler(p any) (err error) {
	// TODO: log the details of the panic, stack trace, etc.
	// For now, just return an internal error with the panic message.
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

// NewUnaryServerInterceptorPanicHandler creates a new gRPC server option for unary interceptor error handling.
func NewUnaryServerInterceptorPanicHandler() (grpc.UnaryServerInterceptor, error) {
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(panicHandler),
	}

	return recovery.UnaryServerInterceptor(opts...), nil
}

// NewStreamServerInterceptorPanicHandler creates a new gRPC server option for unary interceptor error handling.
func NewStreamServerInterceptorPanicHandler() (grpc.StreamServerInterceptor, error) {
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(panicHandler),
	}

	return recovery.StreamServerInterceptor(opts...), nil
}
