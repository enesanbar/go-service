package grpc

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// ServerOptionUnaryInterceptorParams holds the parameters for creating a gRPC server option for unary interceptor.
type ServerOptionUnaryInterceptorParams struct {
	fx.In

	GRPCUnaryServerInterceptors []grpc.UnaryServerInterceptor `group:"grpc-unary-server-interceptors"`
}

// NewServerOptionUnaryInterceptor creates a new gRPC server option for unary interceptor.
func NewServerOptionUnaryInterceptor(p ServerOptionUnaryInterceptorParams) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(p.GRPCUnaryServerInterceptors...)
}
