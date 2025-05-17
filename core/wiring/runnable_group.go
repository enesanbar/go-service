package wiring

import (
	"context"
	"go.uber.org/fx"
)

// RunnableGroup is defined to collect all runnables that may be defined
// For example, HTTP server, Telemetry Server, Health Check Server are all runnables.
type RunnableGroup struct {
	fx.Out

	Runnable Runnable `group:"runnables"`
}

type Runnable interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
