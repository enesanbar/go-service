package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	oteltracing "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/stats/opentelemetry"
)

type GRPCServerOptionOTELParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
}

func NewGRPCServerOptionOTEL(p GRPCServerOptionOTELParams) grpc.ServerOption {
	exporter, err := prometheus.New()
	if err != nil {
		p.Logger.Bg().With(zap.Error(err)).Error("failed to create prometheus exporter")
	}
	// Configure meter provider for metrics
	meterProvider := otelmetric.NewMeterProvider(otelmetric.WithReader(exporter))

	// Configure W3C Trace Context Propagator for traces
	return opentelemetry.ServerOption(opentelemetry.Options{
		MetricsOptions: opentelemetry.MetricsOptions{
			MeterProvider: meterProvider,
		},
		TraceOptions: oteltracing.TraceOptions{
			TracerProvider:    p.TracerProvider,
			TextMapPropagator: p.Propagator,
		}},
	)
}
