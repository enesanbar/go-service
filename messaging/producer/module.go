package producer

import "go.uber.org/fx"

var Module = fx.Module(
	"producer",
	fx.Provide(NewRabbitMQProducer),
)
