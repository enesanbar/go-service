package grpc

import (
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	oteltracing "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/stats/opentelemetry"
)

// ServerOptionOTELParams holds the parameters for creating a gRPC server option for OpenTelemetry tracing and metrics.
type ServerOptionOTELParams struct {
	fx.In

	TracerProvider     *trace.TracerProvider
	Propagator         propagation.TextMapPropagator
	PrometheusExporter *prometheus.Exporter
	MeterProvider      *otelmetric.MeterProvider
}

// NewServerOptionOTEL creates a new gRPC server option for OpenTelemetry tracing and metrics.
func NewServerOptionOTEL(p ServerOptionOTELParams) grpc.ServerOption {
	return opentelemetry.ServerOption(opentelemetry.Options{
		MetricsOptions: opentelemetry.MetricsOptions{
			MeterProvider: p.MeterProvider,
		},
		TraceOptions: oteltracing.TraceOptions{
			TracerProvider:    p.TracerProvider,
			TextMapPropagator: p.Propagator,
		}},
	)
}
