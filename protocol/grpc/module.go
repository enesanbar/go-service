package grpc

import (
	"github.com/enesanbar/go-service/core/service"
	"github.com/enesanbar/go-service/core/wiring"
	"go.uber.org/fx"
)

var module = fx.Module(
	"transport.grpc",
	fx.Provide(
		// provide gRPC server as *Server and wiring.Runnable
		// because gRPC server is needed for registering other services
		// and runnable is needed for starting the server
		NewServer,
		fx.Annotate(
			func(server *Server) *Server {
				return server
			},
			fx.As(new(wiring.Runnable)),
			fx.ResultTags(`group:"runnables"`),
		),
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
		// TODO: The order of interceptors matters. FX adds them randomly. Fix the order.
		AsUnaryServerInterceptor(NewUnaryServerInterceptorProtoValidate),
		AsUnaryServerInterceptor(NewUnaryServerInterceptorErrorHandler),
		AsUnaryServerInterceptor(NewUnaryServerInterceptorPanicHandler),
		AsStreamServerInterceptor(NewStreamServerInterceptorPanicHandler),
		AsServerOption(NewServerOptionUnaryInterceptor),
		AsServerOption(NewServerOptionStreamInterceptor),

		// AsClientOption(NewGRPCClientOptionOTEL), // Experimental
		AsClientOption(NewClientOptionOTELStats),
		// AsClientOption(NewGRPCClientOptionStats),
		AsClientOption(NewClientOptionKeepAliveParams),
		AsClientOption(NewClientOptionCredentials),
		AsClientOption(NewClientOptionCircuitBreaker),
	),
	fx.Invoke(NewHealthCheckHandler), // TODO: Turn this into scheduled task to check health periodically
)

// AsUnaryServerInterceptor is used to annotate a gRPC unary server interceptor for fx.
func AsUnaryServerInterceptor(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-unary-server-interceptors"`),
	)
}

// AsStreamServerInterceptor is used to annotate a gRPC unary server interceptor for fx.
func AsStreamServerInterceptor(p any) any {
	return fx.Annotate(
		p,
		fx.ResultTags(`group:"grpc-stream-server-interceptors"`),
	)
}

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
