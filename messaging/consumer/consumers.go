package consumer

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/wiring"
	"github.com/enesanbar/go-service/messaging/rabbitmq"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

type RabbitMQConsumersParams struct {
	fx.In

	Conf            config.Config
	Logger          log.Factory
	Queues          map[string]*rabbitmq.Queue
	Channels        map[string]*rabbitmq.Channel
	MessageHandlers map[string]MessageHandler
	Propagator      propagation.TextMapPropagator
	TracerProvider  *tracesdk.TracerProvider
}

func RabbitMQConsumerFactory(queueName string, p RabbitMQConsumersParams) (wiring.RunnableGroup, error) {
	runnable, consumer := NewRabbitMQConsumer(RabbitMQConsumersParams{
		Logger:          p.Logger,
		Queues:          p.Queues,
		Channels:        p.Channels,
		MessageHandlers: p.MessageHandlers,
		Propagator:      p.Propagator,
		TracerProvider:  p.TracerProvider,
	})
	consumer.SetQueue(queueName)
	return runnable, nil
}
