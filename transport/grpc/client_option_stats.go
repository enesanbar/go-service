package grpc

import (
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCClientOptionStatsParams struct {
	fx.In

	StatsHandler *StatsHandler
}

func NewGRPCClientOptionStats(p GRPCClientOptionStatsParams) grpc.DialOption {
	// TODO: provide stats on debug mode
	return grpc.WithStatsHandler(p.StatsHandler)
}
