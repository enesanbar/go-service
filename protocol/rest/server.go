package rest

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/wiring"
)

type Server struct {
	logger log.Factory
	server *http.Server
	cfg    *ServerConfig
}

// New creates a pointer to the new instance of the Server
func New(router http.Handler, logger log.Factory, cfg *ServerConfig) (wiring.RunnableGroup, *Server) {
	server := &Server{
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

func (h *Server) Start(ctx context.Context) error {
	h.logger.For(ctx).Infof("starting HTTP Server on %d", h.cfg.Port)
	err := h.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.logger.For(ctx).With(zap.Error(err)).Error("failed to start HTTP server")
		return err
	}
	return nil
}

func (h *Server) Stop(ctx context.Context) error {
	timer := time.AfterFunc(time.Duration(h.cfg.GracefulStopTimeoutSeconds)*time.Second, func() {
		h.logger.For(ctx).Info("http server could not be stopped gracefully, forcing stop")
		err := h.server.Close()
		if err != nil {
			h.logger.For(ctx).With(zap.Error(err)).Error("error forcing http server to stop")
		} else {
			h.logger.For(ctx).Info("http server forced to stop")
		}
	})
	defer timer.Stop()

	h.logger.For(ctx).Info("gracefully stopping HTTP Server")
	if err := h.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error shutting down Server (%w)", err)
	}
	h.logger.For(ctx).Info("HTTP server stopped gracefully")
	return nil
}
