package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

// ClientOptionOTELStatsParams holds the parameters for creating a gRPC dial option for OpenTelemetry stats.
type ClientOptionOTELStatsParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	MeterProvider  *otelmetric.MeterProvider
}

// NewClientOptionOTELStats creates a new gRPC dial option for OpenTelemetry stats.
func NewClientOptionOTELStats(p ClientOptionOTELStatsParams) grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(p.TracerProvider),
		otelgrpc.WithPropagators(p.Propagator),
		otelgrpc.WithMeterProvider(p.MeterProvider),
	))
}
