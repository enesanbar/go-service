package grpc

import (
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	oteltracing "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/stats/opentelemetry"
)

// ClientOptionOTELParams holds the parameters for creating a gRPC dial option for OpenTelemetry tracing and metrics.
type ClientOptionOTELParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	MeterProvider  *otelmetric.MeterProvider
}

// NewClientOptionOTEL creates a new gRPC dial option for OpenTelemetry tracing and metrics.
func NewClientOptionOTEL(p ClientOptionOTELParams) grpc.DialOption {
	return opentelemetry.DialOption(opentelemetry.Options{
		MetricsOptions: opentelemetry.MetricsOptions{
			MeterProvider: p.MeterProvider,
		},
		TraceOptions: oteltracing.TraceOptions{
			TracerProvider:    p.TracerProvider,
			TextMapPropagator: p.Propagator,
		}},
	)
}
