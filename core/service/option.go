package service

import (
	"github.com/enesanbar/go-service/core/cache"
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/healthchecker"
	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/instrumentation/otel"
	"github.com/enesanbar/go-service/core/instrumentation/profiling"
	"github.com/enesanbar/go-service/core/instrumentation/prometheus"
	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/validation"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type Service struct {
	Name         string
	Dependencies fx.Option
	Constructor  interface{}
	InvokeFunc   interface{}
}

type AppConfig struct {
	provides []interface{}
	invokes  []fx.Option
	Options  []fx.Option
	objects  []interface{}
}

type Option func(*AppConfig)

func WithConstructor(constructors ...interface{}) Option {
	return func(cfg *AppConfig) {
		cfg.provides = append(cfg.provides, constructors...)
	}
}

func WithInvoke(funcs ...interface{}) Option {
	return func(cfg *AppConfig) {
		cfg.invokes = append(cfg.invokes, fx.Invoke(funcs...))
	}
}

func WithModules(modules ...fx.Option) Option {
	return func(cfg *AppConfig) {
		cfg.Options = append(cfg.Options, modules...)
	}
}

func WithObject(objects ...interface{}) Option {
	return func(cfg *AppConfig) {
		cfg.objects = append(cfg.objects, objects...)
	}
}

func WithService(service Service) Option {
	return func(cfg *AppConfig) {
		cfg.provides = append(cfg.provides, service.Constructor)
		if service.Dependencies != nil {
			cfg.Options = append(cfg.Options, service.Dependencies)
		}
		if service.InvokeFunc != nil {
			cfg.invokes = append(cfg.invokes, fx.Invoke(service.InvokeFunc))
		}
	}
}

func New(name string, opts ...Option) *fx.App {
	info.ServiceName = name

	cfg := &AppConfig{
		provides: []interface{}{
			log.NewZapLogger,
			log.NewFactory,
		},
		Options: []fx.Option{
			otel.Module,
			prometheus.Module,
			profiling.Module,
			config.Module,
			healthchecker.Module,
			validation.Module,
			cache.Module,
			fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			}),
		},
		invokes: []fx.Option{
			fx.Invoke(bootstrap),
		},
		objects: []interface{}{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return fx.New(
		fx.Provide(cfg.provides...),
		fx.Options(cfg.invokes...),
		fx.Options(cfg.Options...),
		fx.Supply(cfg.objects...),
	)
}
