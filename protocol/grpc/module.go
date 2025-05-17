package grpc

import (
	"github.com/enesanbar/go-service/core/service"
	"go.uber.org/fx"
)

var module = fx.Module(
	"transport.grpc",
	fx.Provide(
		NewServer,
		NewServerConfig,
		NewClientFactory,
		NewRequestLoggerStatsHandler,

		// AsServerOption(NewGRPCServerOptionOTEL), // Experimental
		AsServerOption(NewServerOptionOTELStats),
		// AsServerOption(NewGRPCServerOptionStats),
		AsServerOption(NewServerOptionKeepAliveEnforcementPolicy),
		AsServerOption(NewServerOptionKeepAliveParams),
		AsServerOption(NewServerOptionCredentials),
		AsServerOption(NewServerOptionRequestLoggerStats),
		AsServerOption(NewServerOptionUnaryInterceptorErrorHandler),

		// AsClientOption(NewGRPCClientOptionOTEL), // Experimental
		AsClientOption(NewClientOptionOTELStats),
		// AsClientOption(NewGRPCClientOptionStats),
		AsClientOption(NewClientOptionKeepAliveParams),
		AsClientOption(NewClientOptionCredentials),
		AsClientOption(NewClientOptionCircuitBreaker),
	),
	fx.Invoke(NewHealthCheckHandler), // TODO: Turn this into scheduled task to check health periodically
)

// AsServerOption is used to annotate a gRPC server option for fx.
func AsServerOption(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-server-options"`),
	)
}

// AsClientOption is used to annotate a gRPC client option for fx.
func AsClientOption(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-client-options"`),
	)
}

// Option is used to add gRPC options to the core service.
func Option(options ...fx.Option) service.Option {
	return func(cfg *service.AppConfig) {
		cfg.Options = append(cfg.Options, module)
		cfg.Options = append(cfg.Options, options...)
	}
}
