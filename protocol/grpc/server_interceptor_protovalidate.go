package grpc

import (
	"fmt"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"

	"google.golang.org/grpc"
)

// NewUnaryServerInterceptorProtoValidate creates a new gRPC server option for unary interceptor error handling.
func NewUnaryServerInterceptorProtoValidate() (grpc.UnaryServerInterceptor, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create protovalidate validator: %w", err)
	}

	return protovalidate_middleware.UnaryServerInterceptor(validator), nil
}
