package grpc

import (
	"github.com/enesanbar/go-service/core/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"transport.grpc",
	fx.Provide(
		NewServer,
		NewServerConfig,
		NewClientFactory,
		NewRequestLoggerStatsHandler,

		// AsServerOption(NewGRPCServerOptionOTEL), // Experimental
		AsServerOption(NewGRPCServerOptionOTELStats),
		// AsServerOption(NewGRPCServerOptionStats),
		AsServerOption(NewGRPCServerOptionKeepAliveEnforcementPolicy),
		AsServerOption(NewGRPCServerOptionKeepAliveParams),
		AsServerOption(NewGRPCServerOptionCredentials),
		AsServerOption(NewGRPCServerOptionRequestLoggerStats),
		AsServerOption(NewGRPCServerOptionUnaryInterceptorErrorHandler),

		// AsClientOption(NewGRPCClientOptionOTEL), // Experimental
		AsClientOption(NewGRPCClientOptionOTELStats),
		// AsClientOption(NewGRPCClientOptionStats),
		AsClientOption(NewGRPCClientOptionKeepAliveParams),
		AsClientOption(NewGRPCClientOptionCredentials),
		AsClientOption(NewGRPCClientOptionCircuitBreaker),
	),
	fx.Invoke(NewHealthCheckHandler), // TODO: Turn this into scheduled task to check health periodically
)

func AsServerOption(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-server-options"`),
	)
}

func AsClientOption(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-client-options"`),
	)
}

func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, Module)
		cfg.Options = append(cfg.Options, options...)
	}
}
