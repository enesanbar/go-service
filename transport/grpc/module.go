package grpc

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"transport.grpc",
	fx.Provide(
		New,
		NewConfig,
	),
)
