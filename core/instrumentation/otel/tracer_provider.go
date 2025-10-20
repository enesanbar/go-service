package otel

import (
	"github.com/enesanbar/go-service/core/info"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

type TracerProviderParams struct {
	fx.In

	Exporter    trace.SpanExporter `optional:"true"`
	Environment string             `name:"environment"`
}

func NewTracerProvider(p TracerProviderParams) *trace.TracerProvider {
	opts := []trace.TracerProviderOption{
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(info.ServiceName),
			attribute.String("environment", p.Environment),
		)),
		trace.WithSampler(trace.AlwaysSample()),
	}
	if p.Exporter != nil {
		opts = append(opts, trace.WithBatcher(p.Exporter))
	}

	provider := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(provider)
	return provider
}
