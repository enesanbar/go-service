package grpc

import (
	"github.com/enesanbar/go-service/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCServerOptionOTELStatsParams struct {
	fx.In

	Logger       log.Factory
	Config       *ServerConfig
	StatsHandler *StatsHandler
}

func NewGRPCServerOptionOTELStats(p GRPCServerOptionOTELStatsParams) grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler())
}
