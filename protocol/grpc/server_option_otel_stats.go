package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

// ServerOptionOTELStatsParams holds the parameters for creating a gRPC server option for OpenTelemetry stats.
type ServerOptionOTELStatsParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	MeterProvider  *otelmetric.MeterProvider
}

// NewServerOptionOTELStats creates a new gRPC server option for OpenTelemetry stats.
func NewServerOptionOTELStats(p ServerOptionOTELStatsParams) grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(p.TracerProvider),
		otelgrpc.WithPropagators(p.Propagator),
		otelgrpc.WithMeterProvider(p.MeterProvider),
	))
}
