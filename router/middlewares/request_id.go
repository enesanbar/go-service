package middlewares

import (
	"context"

	"github.com/enesanbar/go-service/utils"

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
			return next(c)
		}
	}
}
