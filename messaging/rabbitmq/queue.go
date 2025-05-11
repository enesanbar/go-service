package rabbitmq

import (
	"github.com/enesanbar/go-service/core/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Queue struct {
	logger  log.Factory
	Queue   *amqp.Queue
	Channel *Channel
	Config  *QueueConfig
}

type QueueParams struct {
	Channel *Channel
	Logger  log.Factory
	Config  *QueueConfig
}

func NewQueue(p QueueParams) (*Queue, error) {
	q, err := p.Channel.Channel.QueueDeclare(
		p.Config.Name,       // name
		p.Config.Durable,    // durable
		p.Config.AutoDelete, // delete when unused
		p.Config.Exclusive,  // exclusive
		p.Config.NoWait,     // no-wait
		nil,                 // arguments
	)

	if err != nil {
		p.Logger.Bg().
			With(zap.String("name", p.Config.Name)).
			With(zap.String("connection", p.Channel.Config.Connection.Config.Host)).
			With(zap.String("channel", p.Channel.Config.Name)).
			With(zap.String("queue", p.Config.Name)).
			With(zap.Error(err)).
			Error("failed to create queue")
		return nil, err
	}

	p.Logger.Bg().
		With(zap.String("name", p.Config.Name)).
		With(zap.String("connection", p.Channel.Config.Connection.Config.Host)).
		With(zap.String("channel", p.Channel.Config.Name)).
		With(zap.String("queue", p.Config.Name)).
		Info("Queue created")

	return &Queue{
		Queue:   &q,
		Channel: p.Channel,
		logger:  p.Logger,
		Config:  p.Config,
	}, nil
}
