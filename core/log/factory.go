package log

import (
	"context"

	"github.com/enesanbar/go-service/core/utils"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	traceIDKey = zap.String("trace_id", "")
	spanIDKey  = zap.String("span_id", "")
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) Factory {
	return Factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
// TODO: investigate alternative logging implementations
func (b Factory) For(ctx context.Context) Logger {
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		logger := spanLogger{span: span, logger: b.logger}

		logger.spanFields = []zapcore.Field{
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
			// zap.Any("context", span.SpanContext()), // debugging purposes
		}

		return logger
	}

	keys := map[string]utils.ContextKey{
		"request_id": utils.ContextKeyRequestID,
	}

	logger := b.Bg()
	for k, v := range keys {
		if value, ok := utils.GetValueFromContext(ctx, v); ok {
			logger = logger.With(zap.String(k, value))
		}
	}

	return logger
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{logger: b.logger.With(fields...)}
}
