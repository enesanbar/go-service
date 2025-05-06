package grpc

import (
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCClientOptionStatsLoggerParams struct {
	fx.In

	StatsHandler *LoggerStatsHandler
}

func NewGRPCClientOptionLoggerStats(p GRPCClientOptionStatsLoggerParams) grpc.DialOption {
	// TODO: provide stats on debug mode
	return grpc.WithStatsHandler(p.StatsHandler)
}
