package service

import (
	"context"
	"reflect"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/wiring"
)

type params struct {
	fx.In

	Runnables       []wiring.Runnable     `group:"runnables"`
	RunnableGroup   [][]wiring.Runnable   `group:"runnable-group"`
	Connections     []wiring.Connection   `group:"connections"`
	ConnectionGroup [][]wiring.Connection `group:"connection-group"`
	Logger          log.Factory
}

// bootstrap defines the fx lifecycle functions OnStart and OnStop
func bootstrap(lc fx.Lifecycle, p params) {
	lc.Append(fx.Hook{
		OnStart: start(p),
		OnStop:  stop(p),
	})
}

// start calls Start method of all types of Runnable and Connection in a separate go routine
func start(p params) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Start all connections
		for _, cg := range p.ConnectionGroup {
			for _, c := range cg {
				p.Connections = append(p.Connections, c)
			}
		}
		for _, connection := range p.Connections {
			if reflect.ValueOf(connection).IsNil() {
				continue
			}
			go func() {
				if err := connection.Start(ctx); err != nil {
					p.Logger.For(ctx).Error("Unable to bootstrap connection", zap.Error(err))
					panic(err)
				}
			}()
		}

		// Start all runnables
		for _, rg := range p.RunnableGroup {
			for _, r := range rg {
				p.Runnables = append(p.Runnables, r)
			}
		}

		for _, r := range p.Runnables {
			if reflect.ValueOf(r).IsNil() {
				continue
			}
			go func(r wiring.Runnable) {
				if err := r.Start(ctx); err != nil {
					p.Logger.For(ctx).Error("Unable to bootstrap runnable", zap.Error(err))
					panic(err)
				}
			}(r)
		}

		return nil
	}
}

// stop calls Stop() method of all types of Runnable and Connection
func stop(p params) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Stop all runnables
		for _, rg := range p.RunnableGroup {
			for _, r := range rg {
				p.Runnables = append(p.Runnables, r)
			}
		}
		wg := sync.WaitGroup{}
		wg.Add(len(p.Runnables))
		for _, server := range p.Runnables {
			go func(server wiring.Runnable) {
				defer wg.Done()
				if reflect.ValueOf(server).IsNil() {
					return
				}
				if err := server.Stop(ctx); err != nil {
					p.Logger.For(ctx).Error("Unable to stop the runnable", zap.Error(err))
				}
			}(server)
		}
		wg.Wait()

		// Stop all connections
		for _, cg := range p.ConnectionGroup {
			for _, c := range cg {
				p.Connections = append(p.Connections, c)
			}
		}
		wg.Add(len(p.Connections))
		for _, connection := range p.Connections {
			go func(connection wiring.Connection) {
				defer wg.Done()
				if reflect.ValueOf(connection).IsNil() {
					return
				}
				if err := connection.Close(ctx); err != nil {
					p.Logger.For(ctx).
						With(zap.String("name", connection.Name())).
						Error("unable to close the connection", zap.Error(err))
				}
			}(connection)
		}
		wg.Wait()
		return nil
	}
}
