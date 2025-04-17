package rabbitmq

import (
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
