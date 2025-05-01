package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCClientOptionOTELStatsParams struct {
	fx.In

	StatsHandler *StatsHandler
}

func NewGRPCClientOptionOTELStats(p GRPCClientOptionOTELStatsParams) grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}
