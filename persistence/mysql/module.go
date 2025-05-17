package mysql

import (
	"github.com/enesanbar/go-service/core/service"
	"github.com/enesanbar/go-service/core/wiring"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"persistence.mysql",
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
)

func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}
