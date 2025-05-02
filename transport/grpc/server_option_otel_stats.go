package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCServerOptionOTELStatsParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	MeterProvider  *otelmetric.MeterProvider
}

func NewGRPCServerOptionOTELStats(p GRPCServerOptionOTELStatsParams) grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(p.TracerProvider),
		otelgrpc.WithPropagators(p.Propagator),
		otelgrpc.WithMeterProvider(p.MeterProvider),
	))
}
