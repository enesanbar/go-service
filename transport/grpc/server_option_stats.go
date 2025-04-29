package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type GRPCServerOptionStatsParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig

	TracerProvider *trace.TracerProvider
	Propagator     propagation.TextMapPropagator
}

func NewGRPCServerOptionStats(p GRPCServerOptionStatsParams) grpc.ServerOption {
	// TODO: provide stats on debug mode
	return grpc.StatsHandler(NewStatsHandler())
}
