package rabbitmq

import (
	"fmt"

	"github.com/enesanbar/go-service/core/messaging/consumer"
	"github.com/enesanbar/go-service/core/messaging/producer"
	"github.com/enesanbar/go-service/core/service"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Conn *Connection `name:"my-connection-1"`
}

var Module = fx.Module(
	"rabbitmq",
	fx.Provide(RabbitMQConnections),
	fx.Provide(RabbitMQChannels),
	fx.Provide(RabbitMQQueues),
	fx.Provide(RabbitMQExchanges),
	fx.Invoke(RabbitMQBindings),
)

func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}

var ProducerModule = fx.Module(
	"messaging.rabbitmq.producer",
	fx.Provide(fx.Annotate(
		NewRabbitMQProducer,
		fx.As(new(producer.Producer)),
	)),
)

type MessageHandlerParams struct {
	fx.In

	Handlers []consumer.MessageHandler `group:"message-handlers"`
}

var ConsumerModule = fx.Module(
	"messaging.rabbitmq.consumer",
	fx.Provide(func(p MessageHandlerParams) map[string]consumer.MessageHandler {
		handlersMap := make(map[string]consumer.MessageHandler)
		for _, handler := range p.Handlers {
			key := fmt.Sprintf("%s-%s", handler.Properties().QueueName, handler.Properties().MessageName)
			handlersMap[key] = handler
		}
		return handlersMap
	}),
)
