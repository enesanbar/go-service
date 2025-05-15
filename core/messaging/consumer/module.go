package consumer

import (
	"fmt"

	"github.com/enesanbar/go-service/core/service"
	"go.uber.org/fx"
)

type MessageHandlerParams struct {
	fx.In

	Handlers []MessageHandler `group:"message-handlers"`
}

var Module = fx.Module(
	"consumer",
	fx.Provide(func(p MessageHandlerParams) map[string]MessageHandler {
		handlersMap := make(map[string]MessageHandler)
		for _, handler := range p.Handlers {
			key := fmt.Sprintf("%s-%s", handler.Properties().QueueName, handler.Properties().MessageName)
			handlersMap[key] = handler
		}
		return handlersMap
	}),
)

func WithConsumer(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}
