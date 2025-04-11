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
	// fx.Invoke(func(p Params) {
	// 	println("Connected to RabbitMQ:", p.Conn.Config.Name)
	// }),
)
