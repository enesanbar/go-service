package grpc

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"transport.grpc",
	fx.Provide(
		New,
		NewConfig,
		AsServerOption(NewGRPCServerOptionOTEL),
		AsServerOption(NewGRPCServerOptionStats),
		AsServerOption(NewGRPCServerOptionKeepAliveEnforcementPolicy),
		AsServerOption(NewGRPCServerOptionKeepAliveParams),
	),
	fx.Invoke(NewHealthCheckHandler), // TODO: Turn this into scheduled task to check health periodically
)

func AsServerOption(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-server-options"`),
	)
}
