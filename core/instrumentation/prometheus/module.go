package prometheus

import (
	"github.com/enesanbar/go-service/core/wiring"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"core/instrumentation/prometheus",
	fx.Provide(
		NewTelemetryServerConfig,
		fx.Annotate(
			NewTelemetryServer,
			fx.As(new(wiring.Runnable)),
			fx.ResultTags(`group:"runnables"`),
		),
	),
)
