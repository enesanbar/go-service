package rabbitmq

import (
	"errors"
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ChannelsParams struct {
	fx.In

	Conf        config.Config
	Logger      log.Factory
	Connections map[string]*Connection `optional:"true"`
}

// Channels return a map of channels configured in the configuration file.
func Channels(p ChannelsParams) (map[string]*Channel, error) {
	if len(p.Connections) == 0 {
		return nil, errors.New("no connections found. please check the connection configuration in your configuration")
	}

	cfg := p.Conf.GetStringMap("rabbitmq.channels")

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
			AppStopSignal: make(chan struct{}),
		}
		channel.connect()
		channels[channelName] = channel
	}

	return channels, nil
}
