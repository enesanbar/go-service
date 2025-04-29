package middlewares

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	AsMiddleware(NewInstanaMiddleware),
	AsMiddleware(NewRequestIDMiddleware),
	AsMiddleware(NewLoggerMiddleware),
	AsMiddleware(NewBodyDumpMiddleware),
	// AsMiddleware(NewMetricsMiddleware),
	AsMiddleware(NewEchoPrometheusMiddleware),
)
