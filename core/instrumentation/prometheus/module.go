package prometheus

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"core/instrumentation/prometheus",
	fx.Provide(
		NewTelemetryServerConfig,
		NewTelemetryServer,
	),
)
