package mongodb

import (
	"github.com/enesanbar/go-service/core/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"persistence.mongodb",
	fx.Provide(NewConnector),
)

func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}
