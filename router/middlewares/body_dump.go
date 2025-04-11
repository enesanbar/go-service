package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// NewBodyDumpMiddleware returns an echo middleware that prints
// request and responses in development mode
func NewBodyDumpMiddleware(p Params) echo.MiddlewareFunc {
	if !p.BaseConfig.IsVerbose() {
		return nil
	}
	return middleware.BodyDump(func(context echo.Context, req []byte, res []byte) {
		p.Logger.
			For(context.Request().Context()).
			With(zap.String("request", string(req))).
			Info("Request Body")

		p.Logger.
			For(context.Request().Context()).
			With(zap.String("response", string(res))).
			Info("Response Body")
	})
}
