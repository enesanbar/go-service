package grpc

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"transport.grpc",
	fx.Provide(
		NewServer,
		NewServerConfig,
		NewClientFactory,
		NewStatsHandler,

		// OTEL Option and OTELStats: may duplicate the spans, use one of them if so
		// AsServerOption(NewGRPCServerOptionOTEL), // Experimental
		AsServerOption(NewGRPCServerOptionOTELStats),
		AsServerOption(NewGRPCServerOptionStats),
		AsServerOption(NewGRPCServerOptionKeepAliveEnforcementPolicy),
		AsServerOption(NewGRPCServerOptionKeepAliveParams),

		// OTEL Option and OTELStats: may duplicate the spans, use one of them if so
		// AsClientOption(NewGRPCClientOptionOTEL), // Experimental
		AsClientOption(NewGRPCClientOptionOTELStats),
		AsClientOption(NewGRPCClientOptionStats),
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
