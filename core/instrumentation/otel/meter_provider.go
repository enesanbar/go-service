package otel

import (
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

type MeterProviderParams struct {
	fx.In

	PrometheusExporter *prometheus.Exporter
}

func NewMeterProvider(p MeterProviderParams) *otelmetric.MeterProvider {
	return otelmetric.NewMeterProvider(
		otelmetric.WithReader(p.PrometheusExporter),
	)
}
