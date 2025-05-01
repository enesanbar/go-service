package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type GRPCClientFactoryParams struct {
	fx.In

	Logger        log.Factory
	ClientOptions []grpc.DialOption `group:"grpc-client-options"`
}

type GRPCClientFactory struct {
	logger log.Factory

	ClientOptions []grpc.DialOption
}

func NewClientFactory(p GRPCClientFactoryParams) (*GRPCClientFactory, error) {

	return &GRPCClientFactory{
		logger:        p.Logger,
		ClientOptions: p.ClientOptions,
	}, nil

}

func (c *GRPCClientFactory) NewClientConn(addr string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Set up OpenTelemetry server interceptor
	c.ClientOptions = append(c.ClientOptions, options...)
	conn, err := grpc.NewClient(addr, c.ClientOptions...)

	if err != nil {
		c.logger.Bg().With(zap.Error(err)).Error("failed to create client")
		return nil, err
	}
	return conn, nil
}
