package wiring

import (
	"go.uber.org/fx"
)

// RunnableGroup is defined to collect all runnables that may be defined
// For example, HTTP server, Telemetry Server, Health Check Server are all runnables.
type RunnableGroup struct {
	fx.Out

	Runnable Runnable `group:"runnables"`
}

type Runnable interface {
	Start() error
	Stop() error
}
