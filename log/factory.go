package log

import (
	"context"

	"github.com/enesanbar/go-service/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
func (b Factory) For(ctx context.Context) Logger {
	// if span := opentracing.SpanFromContext(ctx); span != nil {
	// 	logger := spanLogger{span: span, logger: b.logger}

	// 	if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
	// 		logger.spanFields = []zapcore.Field{
	// 			zap.String("trace_id", jaegerCtx.TraceID().String()),
	// 			zap.String("span_id", jaegerCtx.SpanID().String()),
	// 		}
	// 	}

	// 	return logger
	// }

	keys := map[string]utils.ContextKey{
		"request_id": utils.ContextKeyRequestID,
		"username":   utils.ContextKeyUsername,
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
