package grpc

import (
	"fmt"

	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type GRPCClientFactoryParams struct {
	fx.In

	Logger        log.Factory
	Config        config.Config
	ClientOptions []grpc.DialOption `group:"grpc-client-options"`
}

type GRPCClientFactory struct {
	logger log.Factory
	config config.Config

	ClientOptions []grpc.DialOption
}

func NewClientFactory(p GRPCClientFactoryParams) (*GRPCClientFactory, error) {

	return &GRPCClientFactory{
		logger:        p.Logger,
		config:        p.Config,
		ClientOptions: p.ClientOptions,
	}, nil

}

func (c *GRPCClientFactory) NewClientConn(name string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	addr := c.config.GetString(fmt.Sprintf("client.grpc.%s.address", name))
	c.ClientOptions = append(c.ClientOptions, options...)
	c.ClientOptions = append(c.ClientOptions, grpc.WithDisableServiceConfig())            // Disables service config via TXT record
	c.ClientOptions = append(c.ClientOptions, NewGRPCClientOptionFactory(c.config)(name)) // Adds custom service config
	conn, err := grpc.NewClient(addr, c.ClientOptions...)

	if err != nil {
		c.logger.Bg().With(zap.Error(err)).Error("failed to create client")
		return nil, err
	}
	return conn, nil
}
