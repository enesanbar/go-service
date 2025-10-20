package middlewares

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/enesanbar/go-service/core/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

// NewRequestIDMiddleware returns an echo middleware that
// injects request_id into current request context.
// The value of the incoming request header 'X-Request-ID' is used, if present.
// Otherwise, a new request_id generated.
func NewRequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = random.String(32)
			}
			res.Header().Set(echo.HeaderXRequestID, rid)
			c.Set(echo.HeaderXRequestID, rid)

			ctx := context.WithValue(c.Request().Context(), utils.ContextKeyRequestID, rid)
			requestWithContext := c.Request().WithContext(ctx)
			c.SetRequest(requestWithContext)

			span := trace.SpanFromContext(c.Request().Context())
			if span != nil {
				span.SetAttributes(attribute.String("http.request_id", rid))
			}

			return next(c)
		}
	}
}
