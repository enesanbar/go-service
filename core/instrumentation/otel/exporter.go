package otel

import (
	"context"

	"github.com/enesanbar/go-service/core/osutil"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.uber.org/fx"
)

func NewOTLPExporter() (*otlptrace.Exporter, error) {
	endpointURL := osutil.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT_URL", "http://localhost:4318/v1/traces")

	// further configure the options if needed
	return otlptracehttp.New(
		context.Background(),
		// otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpointURL(endpointURL),
	)
}

func NewStdoutExporter() (*stdout.Exporter, error) {
	return stdout.New(stdout.WithPrettyPrint())
}

func NewZipkinExporter() (*zipkin.Exporter, error) {
	return zipkin.New("http://localhost:9411/api/v2/spans")
}

func NewPrometheusExporter() (*prometheus.Exporter, error) {
	return prometheus.New()
}

func NewExporter() fx.Option {
	exporter := osutil.GetEnv("OTEL_EXPORTER_TYPE", "stdout")
	switch exporter {
	case "otlp":
		return OTLPExporterModule
	case "zipkin":
		return ZipkinExporterModule
	default:
		return StdoutExporterModule
	}
}
