package middlewares

import (
	"github.com/enesanbar/go-service/instrumentation"

	"github.com/labstack/echo/v4"
)

// NewMetricsMiddleware returns an echo middleware that
// saves http RED metrics for prometheus
func NewMetricsMiddleware(service *instrumentation.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if c.Path() == "/metrics" {
				return next(c)
			}

			appMetric := instrumentation.NewHTTP(c.Path(), c.Request().Method)
			if err = next(c); err != nil {
				c.Error(err)
			}
			appMetric.Finished(c)
			go service.SaveMetrics(appMetric)
			return
		}
	}
}
