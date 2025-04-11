package middlewares

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	AsMiddleware(NewInstanaMiddleware),
	AsMiddleware(NewMetricsMiddleware),
	AsMiddleware(NewRequestIDMiddleware),
	AsMiddleware(NewLoggerMiddleware),
	AsMiddleware(NewBodyDumpMiddleware),
)
