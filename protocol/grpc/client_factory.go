package grpc

import (
	"fmt"

	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type ClientFactoryParams struct {
	fx.In

	Logger        log.Factory
	Config        config.Config
	ClientOptions []grpc.DialOption `group:"grpc-client-options"`
}

// ClientFactory is a factory for creating gRPC client connections.
type ClientFactory struct {
	logger log.Factory
	config config.Config

	ClientOptions []grpc.DialOption
}

func NewClientFactory(p ClientFactoryParams) (*ClientFactory, error) {

	return &ClientFactory{
		logger:        p.Logger,
		config:        p.Config,
		ClientOptions: p.ClientOptions,
	}, nil

}

func (c *ClientFactory) NewClientConn(name string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	addr := c.config.GetString(fmt.Sprintf("client.grpc.%s.address", name))
	c.ClientOptions = append(c.ClientOptions, options...)
	c.ClientOptions = append(c.ClientOptions, grpc.WithDisableServiceConfig())                     // Disables service config via TXT record
	c.ClientOptions = append(c.ClientOptions, NewClientOptionServiceConfigFactory(c.config)(name)) // Adds custom service config
	conn, err := grpc.NewClient(addr, c.ClientOptions...)

	if err != nil {
		c.logger.Bg().With(zap.Error(err)).Error("failed to create client")
		return nil, err
	}
	return conn, nil
}
