package rabbitmq

import (
	"fmt"
	"github.com/enesanbar/go-service/core/wiring"

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
	fx.Provide(Connections),
	fx.Provide(
		fx.Annotate(
			func(connections map[string]*Connection) []wiring.Connection {
				result := make([]wiring.Connection, 0, len(connections))
				for _, conn := range connections {
					result = append(result, conn)
				}
				return result
			},
			fx.ResultTags(`group:"connection-group"`),
		),
	),

	fx.Provide(Channels),
	fx.Provide(
		fx.Annotate(
			func(channels map[string]*Channel) []wiring.Connection {
				result := make([]wiring.Connection, 0, len(channels))
				for _, conn := range channels {
					result = append(result, conn)
				}
				return result
			},
			fx.ResultTags(`group:"connection-group"`),
		),
	),

	fx.Provide(Queues),
	fx.Provide(Exchanges),
	fx.Invoke(Bindings),
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

func MapMessageHandlers(p MessageHandlerParams) map[string]consumer.MessageHandler {
	handlersMap := make(map[string]consumer.MessageHandler)
	for _, handler := range p.Handlers {
		key := fmt.Sprintf("%s-%s", handler.Properties().QueueName, handler.Properties().MessageName)
		handlersMap[key] = handler
	}
	return handlersMap
}

var ConsumerModule = fx.Module(
	"messaging.rabbitmq.consumer",
	fx.Provide(MapMessageHandlers),
	fx.Provide(
		fx.Annotate(
			Consumers,
			fx.ResultTags(`group:"runnable-group"`),
		),
	),
)
