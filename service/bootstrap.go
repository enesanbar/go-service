package service

import (
	"context"
	"net/http"
	"reflect"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/wiring"
)

type params struct {
	fx.In

	Runnables   []wiring.Runnable   `group:"runnables"`
	Connections []wiring.Connection `group:"connections"`
	Logger      log.Factory
}

// bootstrap defines the fx lifecycle functions OnStart and OnStop
func bootstrap(lc fx.Lifecycle, p params) {
	lc.Append(fx.Hook{
		OnStart: start(p),
		OnStop:  stop(p),
	})
}

// start calls Start method of all types of Runnable in a separate go routine
func start(p params) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for _, connection := range p.Connections {
			if reflect.ValueOf(connection).IsNil() {
				continue
			}
			go connection.Start(ctx)
		}

		for _, server := range p.Runnables {
			if server == nil {
				continue
			}
			go func(server wiring.Runnable) {
				if err := server.Start(); err != nil && err != http.ErrServerClosed {
					p.Logger.For(ctx).Error("Unable to bootstrap runnable", zap.Error(err))
					panic(err)
				}
			}(server)
		}
		return nil
	}
}

// stop calls Stop() method of all types of Runnable
func stop(p params) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for _, server := range p.Runnables {
			if server == nil {
				continue
			}
			if err := server.Stop(); err != nil {
				p.Logger.For(ctx).Error("Unable to stop the runnable", zap.Error(err))
			}
		}

		for _, connection := range p.Connections {
			if reflect.ValueOf(connection).IsNil() {
				continue
			}
			p.Logger.For(ctx).Info("Closing the connection", zap.String("name", connection.Name()))
			if err := connection.Close(ctx); err != nil {
				p.Logger.For(ctx).Error("unable to close the connection", zap.Error(err))
			}
			p.Logger.For(ctx).Info("Closed the connection", zap.String("name", connection.Name()))
		}
		return nil
	}
}
