package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	oteltracing "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/stats/opentelemetry"
)

type GRPCClientOptionOTELParams struct {
	fx.In

	Logger log.Factory

	TracerProvider     *trace.TracerProvider
	Propagator         propagation.TextMapPropagator
	PrometheusExporter *prometheus.Exporter
}

func NewGRPCClientOptionOTEL(p GRPCClientOptionOTELParams) grpc.DialOption {
	// Configure meter provider for metrics
	meterProvider := otelmetric.NewMeterProvider(otelmetric.WithReader(p.PrometheusExporter))

	// Configure W3C Trace Context Propagator for traces
	return opentelemetry.DialOption(opentelemetry.Options{
		MetricsOptions: opentelemetry.MetricsOptions{
			MeterProvider: meterProvider,
		},
		TraceOptions: oteltracing.TraceOptions{
			TracerProvider:    p.TracerProvider,
			TextMapPropagator: p.Propagator,
		}},
	)
}
