package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCServerOptionLoggerStatsParams struct {
	fx.In

	Logger                    log.Factory
	Config                    *ServerConfig
	RequestLoggerStatsHandler *RequestLoggerStatsHandler
}

func NewGRPCServerOptionRequestLoggerStats(p GRPCServerOptionLoggerStatsParams) grpc.ServerOption {
	return grpc.StatsHandler(p.RequestLoggerStatsHandler)
}
