package router

import (
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/wiring"

	"github.com/enesanbar/go-service/core/log"
)

type ProfileServer struct {
	Router     *http.ServeMux
	logger     log.Factory
	Server     io.Closer
	BaseConfig *config.Base
}

func NewProfileServer(p EchoParams) (wiring.RunnableGroup, *ProfileServer) {
	server := &ProfileServer{
		logger: p.Logger,
	}

	if os.Getenv("PPROF") == "true" || p.BaseConfig.IsVerbose() {
		return wiring.RunnableGroup{Runnable: server}, server
	}

	return wiring.RunnableGroup{}, nil
}

func (ts *ProfileServer) Start() error {
	ts.logger.Bg().Info("starting profiling server on 6060...")
	return http.ListenAndServe(":6060", nil)
}

func (ts *ProfileServer) Stop() error {
	ts.logger.Bg().Info("gracefully stopping profiling server")
	return nil
}
