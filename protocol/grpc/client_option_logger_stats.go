package grpc

import (
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

// ClientOptionStatsLoggerParams holds the parameters for creating a gRPC dial option for the request logger stats handler.
type ClientOptionStatsLoggerParams struct {
	fx.In

	RequestLoggerStatsHandler *RequestLoggerStatsHandler
}

// NewClientOptionLoggerStats creates a new gRPC dial option for the request logger stats handler.
func NewClientOptionLoggerStats(p ClientOptionStatsLoggerParams) grpc.DialOption {
	return grpc.WithStatsHandler(p.RequestLoggerStatsHandler)
}
