package otel

import (
	"github.com/enesanbar/go-service/router/middlewares"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var Module = fx.Options(
	factories,
)

var OTLPExporterModule = fx.Provide(
	fx.Annotate(
		NewOTLPExporter,
		fx.As(new(trace.SpanExporter)),
	),
)

var StdoutExporterModule = fx.Provide(
	fx.Annotate(
		NewStdoutExporter,
		fx.As(new(trace.SpanExporter)),
	),
)

var ZipkinExporterModule = fx.Provide(
	fx.Annotate(
		NewZipkinExporter,
		fx.As(new(trace.SpanExporter)),
	),
)

var factories = fx.Options(
	NewExporter(),
	fx.Provide(
		NewTracerProvider,
		NewPropagator,
		middlewares.AsMiddleware(NewOtelMiddleware),
	),
)
