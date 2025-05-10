package service

import (
	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/healthchecker"
	"github.com/enesanbar/go-service/info"
	"github.com/enesanbar/go-service/instrumentation"
	"github.com/enesanbar/go-service/instrumentation/otel"
	"github.com/enesanbar/go-service/log"
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
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

// Service is a dependency that is required by multiple modules of the
// application, such as logging.
type Service struct {
	Name         string
	Dependencies fx.Option
	Constructor  interface{}
	InvokeFunc   interface{}
}

// New creates a new instance of the Builder type
func New(name string) Builder {
	info.ServiceName = name

	return Builder{
		Provide: []interface{}{
			instrumentation.NewTelemetryServer,
			instrumentation.NewTelemetryServerConfig,
			log.NewZapLogger,
			log.NewFactory,
		},
		Invoke: []fx.Option{
			fx.Invoke(bootstrap),
		},
		Options: []fx.Option{
			validation.Module,
			otel.Module,
			config.Module,
			healthchecker.Module,
			fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			}),
			rabbitmq.Module,
			producer.Module,
			mysql.Module,
		},
		Objects: []interface{}{},
	}
}

// Builder is a type that will build an *fx.App from dependencies, it stores the command
// that started the app, along with all the dependencies required to bootstrap the app.
type Builder struct {
	Provide []interface{}
	Invoke  []fx.Option
	Options []fx.Option
	Objects []interface{}
}

// WithObject is analogous to the fx.Supply option,
// it allows you to provide multiple instantiated object
func (b Builder) WithObject(constructors ...interface{}) Builder {
	b.Provide = append(b.Provide, constructors...)
	return b
}

// WithConstructor is analogous to the fx.Provide option,
// it allows you to provide multiple constructors
func (b Builder) WithConstructor(constructors ...interface{}) Builder {
	b.Provide = append(b.Provide, constructors...)
	return b
}

// WithService allows you to register a Service dependency, this should be done within
// an `invoke()` function, so that the flags from the
func (b Builder) WithService(service Service) Builder {
	b = b.WithConstructor(service.Constructor)
	if service.Dependencies != nil {
		b.Options = append(b.Options, service.Dependencies)
	}
	if service.InvokeFunc == nil {
		return b
	}
	return b.WithInvoke(service.InvokeFunc)
}

func (b Builder) WithRestAdapter(options ...fx.Option) Builder {
	b.Options = append(b.Options, http.Module, router.Module)
	b.Options = append(b.Options, options...)
	return b
}

func (b Builder) WithGRPCAdapter(options ...fx.Option) Builder {
	b.Options = append(b.Options, grpc.Module)
	b.Options = append(b.Options, options...)
	return b
}

func (b Builder) WithConsumer(options ...fx.Option) Builder {
	b.Options = append(b.Options, consumer.Module)
	b.Options = append(b.Options, options...)
	return b
}

// WithModules adds a fx.Option into the list of dependencies to be built,
// used for adding application modules
func (b Builder) WithModules(modules ...fx.Option) Builder {
	b.Options = append(b.Options, modules...)
	return b
}

// WithInvoke will add a function as a function to be invoked
func (b Builder) WithInvoke(funcs ...interface{}) Builder {
	b.Invoke = append(b.Invoke, fx.Invoke(funcs...))
	return b
}

// Build will produce a new instance of the *fx.App from the variables of the builder
func (b Builder) Build() *fx.App {
	return fx.New(
		fx.Provide(b.Provide...),
		fx.Options(b.Invoke...),
		fx.Options(b.Options...),
		fx.Supply(b.Objects...),
	)
}

// BuildTest will create a new instance of *fxtest.App from the contained dependencies
func (b Builder) BuildTest(tb fxtest.TB) *fxtest.App {
	return fxtest.New(
		tb,
		fx.Provide(b.Provide...),
		fx.Options(b.Invoke...),
		fx.Options(b.Options...),
		fx.Supply(b.Objects...),
	)
}
