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

type GRPCClientOptionOTELParams struct {
	fx.In

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
	MeterProvider  *otelmetric.MeterProvider
}

func NewGRPCClientOptionOTEL(p GRPCClientOptionOTELParams) grpc.DialOption {
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
