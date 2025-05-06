package consumer

import (
	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/messaging/rabbitmq"
	"github.com/enesanbar/go-service/wiring"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

type RabbitMQConsumersParams struct {
	fx.In

	QueueName       string
	Conf            config.Config
	Logger          log.Factory
	Queues          map[string]*rabbitmq.Queue
	Channels        map[string]*rabbitmq.Channel
	MessageHandlers map[string]MessageHandler
	Propagator      propagation.TextMapPropagator
	TracerProvider  *tracesdk.TracerProvider
}

func RabbitMQConsumerFactory(p RabbitMQConsumersParams) (wiring.RunnableGroup, error) {
	runnable, consumer := NewRabbitMQConsumer(RabbitMQConsumerParams{
		Logger:          p.Logger,
		Queues:          p.Queues,
		Channels:        p.Channels,
		MessageHandlers: p.MessageHandlers,
		Propagator:      p.Propagator,
		TracerProvider:  p.TracerProvider,
	})
	consumer.SetQueue(p.QueueName)
	return runnable, nil
}
