package inmemory

import (
	"github.com/enesanbar/go-service/core/cache"
	"github.com/enesanbar/go-service/core/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"cache.inmemory",
	fx.Provide(
		NewConfig,
		fx.Annotate(
			NewInMemoryCache,
			fx.As(new(cache.Cache)),
		),
	),
)

func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}
