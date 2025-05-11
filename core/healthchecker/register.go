package healthchecker

import "go.uber.org/fx"

// AsHealthCheckerProbe is a function that provides the constructor function to
// the fx container under the group "health-checker-probes".
// the provided function must implement the [HealthCheckerProbe] interface.
func AsHealthCheckerProbe(p any) any {
	return fx.Annotate(
		p,
		fx.As(new(HealthCheckerProbe)),
		fx.ResultTags(`group:"health-checker-probes"`),
	)
}
