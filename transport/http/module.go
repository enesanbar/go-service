package http

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"transport.http",
	fx.Provide(
		New,
		NewConfig,
	),
)
