package rabbitmq

import (
	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RabbitMQExchangesParams struct {
	fx.In

	Conf     config.Config
	Logger   log.Factory
	Channels map[string]*Channel `optional:"true"`
}

func RabbitMQExchanges(p RabbitMQExchangesParams) (map[string]*Exchange, error) {
	if len(p.Channels) == 0 {
		return nil, nil
	}
	cfg := p.Conf.GetStringMap("datasources.rabbitmq.exchanges")

	exchanges := make(map[string]*Exchange)

	for exchangeName, v := range cfg {
		channelName := v.(map[string]interface{})[PropertyChannel].(string)
		channel, ok := p.Channels[channelName]
		if !ok {
			p.Logger.Bg().
				With(zap.String("exchange", exchangeName)).
				With(zap.String("channel", channelName)).
				Error("channel not found for exchange. please check the channel configuration in your configuration")
			continue
		}

		Exchange, err := NewExchange(ExchangeParams{
			Channel: channel,
			Logger:  p.Logger,
			Config: &ExchangeConfig{
				Name:       exchangeName,
				Channel:    channel,
				Type:       v.(map[string]interface{})[PropertyType].(string),
				Durable:    v.(map[string]interface{})[PropertyDurable].(bool),
				AutoDelete: v.(map[string]interface{})[PropertyAutoDelete].(bool),
				Internal:   v.(map[string]interface{})[PropertyExclusive].(bool),
				NoWait:     v.(map[string]interface{})[PropertyNoWait].(bool),
			},
		})

		if err != nil {
			p.Logger.Bg().
				With(zap.String("exchange", exchangeName)).
				With(zap.String("channel", channelName)).
				With(zap.Error(err)).
				Error("failed to create Exchange")
			panic(err)
		}

		exchanges[exchangeName] = Exchange
		p.Logger.Bg().
			With(zap.String("exchange", exchangeName)).
			With(zap.String("channel", channelName)).
			Info("Exchange created")
	}

	return exchanges, nil
}
