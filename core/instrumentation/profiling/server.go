package profiling

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/enesanbar/go-service/core/wiring"

	"github.com/enesanbar/go-service/core/log"
)

type ProfileServer struct {
	logger log.Factory
}

func NewProfileServer(log log.Factory) (wiring.RunnableGroup, *ProfileServer) {
	server := &ProfileServer{
		logger: log,
	}

	if os.Getenv("PPROF") == "true" {
		return wiring.RunnableGroup{Runnable: server}, server
	}

	return wiring.RunnableGroup{}, nil
}

func (ts *ProfileServer) Start(ctx context.Context) error {
	ts.logger.For(ctx).Info("starting profiling server on 6060...")
	return http.ListenAndServe(":6060", nil)
}

func (ts *ProfileServer) Stop(ctx context.Context) error {
	ts.logger.For(ctx).Info("stopping profiling server")
	return nil
}
