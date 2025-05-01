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
		AsServerOption(NewGRPCServerOptionOTEL),
		AsServerOption(NewGRPCServerOptionOTELStats),
		AsServerOption(NewGRPCServerOptionStats),
		AsServerOption(NewGRPCServerOptionKeepAliveEnforcementPolicy),
		AsServerOption(NewGRPCServerOptionKeepAliveParams),
		AsClientOption(NewGRPCServerOptionStats),
		AsClientOption(NewGRPCClientOptionOTEL),
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
