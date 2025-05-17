package grpc

import (
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

// ServerOptionLoggerStatsParams holds the parameters for creating a gRPC server option for the request logger stats handler.
type ServerOptionLoggerStatsParams struct {
	fx.In

	Logger                    log.Factory
	Config                    *ServerConfig
	RequestLoggerStatsHandler *RequestLoggerStatsHandler
}

// NewServerOptionRequestLoggerStats creates a new gRPC server option for request logger stats.
func NewServerOptionRequestLoggerStats(p ServerOptionLoggerStatsParams) grpc.ServerOption {
	return grpc.StatsHandler(p.RequestLoggerStatsHandler)
}
