package grpc

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// ServerOptionStreamInterceptorParams holds the parameters for creating a gRPC server option for stream interceptor.
type ServerOptionStreamInterceptorParams struct {
	fx.In

	GRPCStreamServerInterceptors []grpc.StreamServerInterceptor `group:"grpc-stream-server-interceptors"`
}

// NewServerOptionStreamInterceptor creates a new gRPC server option for unary interceptor.
func NewServerOptionStreamInterceptor(p ServerOptionStreamInterceptorParams) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(p.GRPCStreamServerInterceptors...)
}
