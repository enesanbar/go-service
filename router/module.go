package router

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/enesanbar/go-service/instrumentation"
	"github.com/enesanbar/go-service/router/middlewares"
)

var Module = fx.Options(
	factories,
	middlewares.Module,
	interfaceTypes,
)

var factories = fx.Provide(
	NewBaseHandler,
	NewHealthCheckHandler,
	NewEchoRouter,
	NewProfileServer,
	instrumentation.NewPrometheusService,
)

var interfaceTypes = fx.Provide(
	func(echoRouter *EchoServer) http.Handler { return echoRouter },
)
