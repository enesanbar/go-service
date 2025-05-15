package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCClientOptionOTELStatsParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	MeterProvider  *otelmetric.MeterProvider
}

func NewGRPCClientOptionOTELStats(p GRPCClientOptionOTELStatsParams) grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(p.TracerProvider),
		otelgrpc.WithPropagators(p.Propagator),
		otelgrpc.WithMeterProvider(p.MeterProvider),
	))
}
