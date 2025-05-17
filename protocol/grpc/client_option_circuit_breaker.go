package grpc

import (
	"context"
	"fmt"

	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/log"
	"github.com/sony/gobreaker"
	"go.uber.org/fx"

	"google.golang.org/grpc"
)

type ClientOptionCircuitBreakerParams struct {
	fx.In

	Logger log.Factory
	Config *ServerConfig
}

func NewClientOptionCircuitBreaker(p ClientOptionCircuitBreakerParams) grpc.DialOption {
	// TODO: configure the circuit breaker settings from the config per service, or per method. use the default config if none is provided
	// map[string]*CircuitBreaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        info.ServiceName,
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.1
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			p.Logger.Bg().Info(fmt.Sprintf("Circuit Breaker: %s, changed from %v, to %v", name, from, to))
		},
	})

	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, cbErr := cb.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return nil, err
			}

			return nil, nil

		})
		return cbErr
	}
	return grpc.WithUnaryInterceptor(interceptor)
}
