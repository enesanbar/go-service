package grpc

import (
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type GRPCClientOptionStatsLoggerParams struct {
	fx.In

	RequestLoggerStatsHandler *RequestLoggerStatsHandler
}

func NewGRPCClientOptionLoggerStats(p GRPCClientOptionStatsLoggerParams) grpc.DialOption {
	return grpc.WithStatsHandler(p.RequestLoggerStatsHandler)
}
