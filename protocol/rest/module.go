package rest

import (
	"github.com/enesanbar/go-service/core/service"
	"github.com/enesanbar/go-service/protocol/rest/router"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"transport.http",
	fx.Provide(
		New,
		NewConfig,
	),
	fx.Options(router.Module),
)

func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}
