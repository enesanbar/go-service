package middlewares

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

// NewLoggerMiddleware returns an echo middleware that prints
// incoming requests to stdout
func NewEchoPrometheusMiddleware(p Params) echo.MiddlewareFunc {
	return echoprometheus.NewMiddleware("http")
}
