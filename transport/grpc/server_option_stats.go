package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCServerOptionStatsParams struct {
	fx.In

	Logger       log.Factory
	Config       *ServerConfig
	StatsHandler *StatsHandler
}

func NewGRPCServerOptionStats(p GRPCServerOptionStatsParams) grpc.ServerOption {
	// TODO: provide stats on debug mode
	return grpc.StatsHandler(p.StatsHandler)
}
