package healthchecker

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"health-checker",
	fx.Provide(
		NewDefaultHealthCheckerFactory,
		NewFxHealthChecker,
	),
)

type FxHealthCheckerParam struct {
	fx.In

	Factory HealthCheckerFactory
	Probes  []HealthCheckerProbe `group:"health-checker-probes"`
}

func NewFxHealthChecker(p FxHealthCheckerParam) (*HealthChecker, error) {
	return p.Factory.Create(WithProbes(p.Probes...))
}
