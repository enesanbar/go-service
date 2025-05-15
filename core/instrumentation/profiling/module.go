package profiling

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"core/instrumentation/profiling",
	fx.Provide(
		NewProfileServer,
	),
)
