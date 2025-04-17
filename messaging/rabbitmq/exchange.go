package rabbitmq

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/zap"
)

type Exchange struct {
	logger  log.Factory
	Channel *Channel
	Config  *ExchangeConfig
}

type ExchangeParams struct {
	Channel *Channel
	Logger  log.Factory
	Config  *ExchangeConfig
}

func NewExchange(p ExchangeParams) (*Exchange, error) {
	err := p.Channel.Channel.ExchangeDeclare(
		p.Config.Name,       // name
		p.Config.Type,       // kind
		p.Config.Durable,    // durable
		p.Config.AutoDelete, // delete when unused
		p.Config.Internal,   // exclusive
		p.Config.NoWait,     // no-wait
		nil,                 // arguments
	)

	if err != nil {
		p.Logger.Bg().
			With(zap.String("name", p.Config.Name)).
			With(zap.String("connection", p.Channel.Config.Connection.Config.Host)).
			With(zap.String("channel", p.Channel.Config.Name)).
			With(zap.String("exchange", p.Config.Name)).
			With(zap.Error(err)).
			Error("failed to create Exchange")
		return nil, err
	}

	p.Logger.Bg().
		With(zap.String("name", p.Config.Name)).
		With(zap.String("connection", p.Channel.Config.Connection.Config.Host)).
		With(zap.String("channel", p.Channel.Config.Name)).
		With(zap.String("exchange", p.Config.Name)).
		Info("Exchange created")

	return &Exchange{
		Channel: p.Channel,
		logger:  p.Logger,
		Config:  p.Config,
	}, nil
}
