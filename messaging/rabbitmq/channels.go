package rabbitmq

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RabbitMQChannelsParams struct {
	fx.In

	Conf        config.Config
	Logger      log.Factory
	Connections map[string]*Connection `optional:"true"`
}

func RabbitMQChannels(p RabbitMQChannelsParams) (map[string]*Channel, error) {
	if len(p.Connections) == 0 {
		return nil, nil
	}

	cfg := p.Conf.GetStringMap("datasources.rabbitmq.channels")

	channels := make(map[string]*Channel)
	for channelName, v := range cfg {
		connectionName := v.(map[string]interface{})["connection"].(string)
		connection, ok := p.Connections[connectionName]
		if !ok {
			p.Logger.Bg().
				With(zap.String("channel", channelName)).
				With(zap.String("connection", connectionName)).
				Error("connection not found for channel. please check the connection configuration in your configuration")
			continue
		}
		channel := &Channel{
			logger: p.Logger,
			Config: &ChannelConfig{
				Name:       channelName,
				Connection: connection,
			},
		}
		channel.connect()
		channels[channelName] = channel
	}

	return channels, nil
}
