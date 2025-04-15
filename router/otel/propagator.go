package otel

import (
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
)

type PropagatorParams struct {
	fx.In
}

func NewPropagator(p PropagatorParams) propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
