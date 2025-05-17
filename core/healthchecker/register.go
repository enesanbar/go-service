package healthchecker

import "go.uber.org/fx"

// AsHealthCheckerProbe is a function that provides the constructor function to
// the fx container under the group "health-checker-probes".
// the provided function must implement the [Probe] interface.
func AsHealthCheckerProbe(p any) any {
	return fx.Annotate(
		p,
		fx.As(new(Probe)),
		fx.ResultTags(`group:"health-checker-probes"`),
	)
}
