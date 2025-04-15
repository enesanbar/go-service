package otel

import (
	"github.com/enesanbar/go-service/info"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

type OtelMiddlewareParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
}

// NewOtelMiddleware returns an echo middleware for OpenTelemetry
func NewOtelMiddleware(p OtelMiddlewareParams) echo.MiddlewareFunc {
	return otelecho.Middleware(
		info.ServiceName,
		otelecho.WithTracerProvider(p.TracerProvider),
		otelecho.WithPropagators(p.Propagator),
	)
}
