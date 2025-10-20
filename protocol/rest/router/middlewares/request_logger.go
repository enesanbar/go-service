package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// NewLoggerMiddleware returns an echo middleware that prints
// incoming requests to stdout
func NewLoggerMiddleware(p Params) echo.MiddlewareFunc {
	if !p.BaseConfig.IsVerbose() {
		return nil
	}
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			p.Logger.For(c.Request().Context()).Info("",
				zap.Duration("latency", v.Latency),
				zap.String("protocol", v.Protocol),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("URI", v.URI),
				zap.String("uri_path", v.URIPath),
				zap.String("route_path", v.RoutePath),
				zap.String("request_id", v.RequestID),
				zap.String("referer", v.Referer),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", v.Status),
				zap.Error(v.Error),
				zap.String("content_length", v.ContentLength),
				zap.Int64("response_size", v.ResponseSize),
			)

			return nil
		},
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
	})
}
