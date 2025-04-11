package middlewares

import (
	instana "github.com/instana/go-sensor"
	"github.com/instana/go-sensor/instrumentation/instaecho"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type InstanaParams struct {
	fx.In

	Sensor *instana.Sensor `optional:"true"`
}

// NewInstanaMiddleware returns an echo middleware that
// adds tracing context and handles entry span.
func NewInstanaMiddleware(p InstanaParams) echo.MiddlewareFunc {
	if p.Sensor == nil {
		return nil
	}
	return instaecho.Middleware(p.Sensor)
}
