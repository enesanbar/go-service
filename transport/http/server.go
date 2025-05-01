package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/wiring"
)

type HttpServer struct {
	logger log.Factory
	server *http.Server
	cfg    *ServerConfig
}

// New creates a pointer to the new instance of the HttpServer
func New(router http.Handler, logger log.Factory, cfg *ServerConfig) (wiring.RunnableGroup, *HttpServer) {
	server := &HttpServer{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		},
		logger: logger,
		cfg:    cfg,
	}
	return wiring.RunnableGroup{
		Runnable: server,
	}, server
}

func (h HttpServer) Start() error {
	h.logger.Bg().Infof("starting HTTP Server on %d", h.cfg.Port)
	return h.server.ListenAndServe()
}

func (h HttpServer) Stop() error {
	timer := time.AfterFunc(time.Duration(h.cfg.GracefulStopTimeoutSeconds)*time.Second, func() {
		h.logger.Bg().Info("http server could not be stopped gracefully, forcing stop")
		h.server.Close()
		h.logger.Bg().Info("http server forced to stop")
	})
	defer timer.Stop()

	h.logger.Bg().Info("gracefully stopping HTTP Server")
	if err := h.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error shutting down Server (%w)", err)
	}
	h.logger.Bg().Info("HTTP server stopped gracefully")
	return nil
}
