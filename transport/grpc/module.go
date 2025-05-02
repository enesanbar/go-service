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

		// AsServerOption(NewGRPCServerOptionOTEL), // Experimental
		AsServerOption(NewGRPCServerOptionOTELStats),
		// AsServerOption(NewGRPCServerOptionStats),
		AsServerOption(NewGRPCServerOptionKeepAliveEnforcementPolicy),
		AsServerOption(NewGRPCServerOptionKeepAliveParams),
		AsServerOption(NewGRPCServerOptionCredentials),

		AsClientOption(NewGRPCClientOptionOTEL), // Experimental
		// AsClientOption(NewGRPCClientOptionOTELStats),
		// AsClientOption(NewGRPCClientOptionStats),
		AsClientOption(NewGRPCClientOptionKeepAliveParams),
		AsClientOption(NewGRPCClientOptionCredentials),
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
