package otel

import (
	"github.com/enesanbar/go-service/info"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

type TracerProviderParams struct {
	fx.In

	Exporter trace.SpanExporter
}

func NewTracerProvider(p TracerProviderParams) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(p.Exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(info.ServiceName),
		)),
		trace.WithSampler(trace.AlwaysSample()),
	)
}
