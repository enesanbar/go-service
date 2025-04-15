package router

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/enesanbar/go-service/instrumentation"
	"github.com/enesanbar/go-service/router/middlewares"
	"github.com/enesanbar/go-service/router/otel"
)

var Module = fx.Options(
	factories,
	middlewares.Module,
	otel.Module,
	interfaceTypes,
)

var factories = fx.Provide(
	NewBaseHandler,
	NewHealthCheckHandler,
	NewEchoRouter,
	NewTelemetryServer,
	NewProfileServer,
	instrumentation.NewPrometheusService,
)

var interfaceTypes = fx.Provide(
	func(echoRouter *EchoServer) http.Handler { return echoRouter },
)
