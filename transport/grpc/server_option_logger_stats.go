package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCServerOptionLoggerStatsParams struct {
	fx.In

	Logger       log.Factory
	Config       *ServerConfig
	StatsHandler *LoggerStatsHandler
}

func NewGRPCServerOptionLoggerStats(p GRPCServerOptionLoggerStatsParams) grpc.ServerOption {
	// TODO: provide stats on debug mode
	return grpc.StatsHandler(p.StatsHandler)
}
