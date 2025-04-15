package router

import (
	"net/http"

	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

// ApplierFunc is a function type that allows routes to be applied to the main router
type ApplierFunc func(router *echo.Group)

// RouteConfig is a group of routes to route to based on a path
type RouteConfig struct {
	Path   string
	Router ApplierFunc
}

type EchoServer struct {
	Router        *echo.Echo
	ContextRouter *echo.Group
	logger        log.Factory
	config        config.Config
}

type EchoParams struct {
	fx.In

	Middlewares          []echo.MiddlewareFunc `group:"middlewares"`
	Routes               []RouteConfig         `group:"routes"`
	Logger               log.Factory
	Config               config.Config
	BaseConfig           *config.Base
	HealthCheckerHandler *HealthCheckHandler
}

func NewEchoRouter(p EchoParams) *EchoServer {
	e := echo.New()

	// apply middlewares
	for _, middleware := range p.Middlewares {
		if middleware != nil {
			e.Use(middleware)
		}
	}

	contextPath := p.Config.GetString("server.context-path")
	contextRouter := e.Group(contextPath)

	healthCheckPath := p.Config.GetString("server.healthcheck-path")
	if healthCheckPath == "" {
		healthCheckPath = "/health"
	}
	contextRouter.GET(healthCheckPath, p.HealthCheckerHandler.Handle)

	// apply routes
	for _, route := range p.Routes {
		group := contextRouter.Group(route.Path)
		route.Router(group)
	}

	// setup swagger route
	contextRouter.GET("/api/swagger/*", echoSwagger.WrapHandler).Name = "swagger"
	contextRouter.GET("/api", func(context echo.Context) error {
		return context.Redirect(
			http.StatusMovedPermanently,
			contextPath+"/api/swagger/index.html")
	})

	return &EchoServer{
		Router:        e,
		ContextRouter: contextRouter,
		logger:        p.Logger,
		config:        p.Config,
	}
}

func (es EchoServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	es.Router.ServeHTTP(writer, request)
}
