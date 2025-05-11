package service

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/healthchecker"
	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/instrumentation/otel"
	"github.com/enesanbar/go-service/core/instrumentation/prometheus"
	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/messaging/consumer"
	"github.com/enesanbar/go-service/messaging/producer"
	"github.com/enesanbar/go-service/messaging/rabbitmq"
	"github.com/enesanbar/go-service/persistance/mysql"
	"github.com/enesanbar/go-service/router"
	"github.com/enesanbar/go-service/transport/grpc"
	"github.com/enesanbar/go-service/transport/http"
	"github.com/enesanbar/go-service/validation"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type AppConfig struct {
	provides []interface{}
	invokes  []fx.Option
	options  []fx.Option
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
		cfg.options = append(cfg.options, modules...)
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
			cfg.options = append(cfg.options, service.Dependencies)
		}
		if service.InvokeFunc != nil {
			cfg.invokes = append(cfg.invokes, fx.Invoke(service.InvokeFunc))
		}
	}
}

func WithGRPCAdapter(options ...fx.Option) Option {
	return func(cfg *AppConfig) {
		cfg.options = append(cfg.options, grpc.Module)
		cfg.options = append(cfg.options, options...)
	}
}

func WithRestAdapter(options ...fx.Option) Option {
	return func(cfg *AppConfig) {
		cfg.options = append(cfg.options, http.Module, router.Module)
		cfg.options = append(cfg.options, options...)
	}
}

func WithConsumer(options ...fx.Option) Option {
	return func(cfg *AppConfig) {
		cfg.options = append(cfg.options, consumer.Module)
		cfg.options = append(cfg.options, options...)
	}
}

func WithMySQL(options ...fx.Option) Option {
	return func(cfg *AppConfig) {
		cfg.options = append(cfg.options, mysql.Module)
		cfg.options = append(cfg.options, options...)
	}
}

func WithRabbitMQ(options ...fx.Option) Option {
	return func(cfg *AppConfig) {
		cfg.options = append(cfg.options, rabbitmq.Module, producer.Module)
		cfg.options = append(cfg.options, options...)
	}
}

func NewApp(name string, opts ...Option) *fx.App {
	info.ServiceName = name

	cfg := &AppConfig{
		provides: []interface{}{
			log.NewZapLogger,
			log.NewFactory,
		},
		options: []fx.Option{
			otel.Module,
			prometheus.Module,
			config.Module,
			healthchecker.Module,
			fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			}),
			validation.Module,
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
		fx.Options(cfg.options...),
		fx.Supply(cfg.objects...),
	)
}
