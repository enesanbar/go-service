package cache

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"core.cache",
	fx.Provide(
		NewInstrumentor,
	),
)
