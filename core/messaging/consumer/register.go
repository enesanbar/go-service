package consumer

import "go.uber.org/fx"

type MessageProperties struct {
	QueueName   string
	MessageName string
}

func AsMessageHandler(p any) any {
	return fx.Annotate(
		p,
		fx.As(new(MessageHandler)),
		fx.ResultTags(`group:"message-handlers"`),
	)
}
