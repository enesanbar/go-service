package prometheus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/wiring"

	"github.com/enesanbar/go-service/core/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TelemetryServer struct {
	Router     *http.ServeMux
	logger     log.Factory
	Server     io.Closer
	BaseConfig *config.Base
	cfg        *TelemetryServerConfig
}

func NewTelemetryServer(
	logger log.Factory,
	baseConfig *config.Base,
	telemetryConfig *TelemetryServerConfig,
) (wiring.RunnableGroup, *TelemetryServer) {
	telemetryRouter := http.NewServeMux()

	telemetryRouter.Handle("/metrics", promhttp.Handler())

	server := &TelemetryServer{
		Router:     telemetryRouter,
		logger:     logger,
		BaseConfig: baseConfig,
		cfg:        telemetryConfig,
	}
	return wiring.RunnableGroup{Runnable: server}, server
}

func (ts *TelemetryServer) Start(ctx context.Context) error {
	logger := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ts.BaseConfig.IsVerbose() {
				ts.logger.For(r.Context()).Info("",
					zap.String("protocol", r.Proto),
					zap.String("remote_ip", r.RemoteAddr),
					zap.String("host", r.Host),
					zap.String("method", r.Method),
					zap.String("URI", r.RequestURI),
					zap.String("referer", r.Referer()),
					zap.String("user_agent", r.UserAgent()),
					zap.Int64("content_length", r.ContentLength),
				)
			}
			handler.ServeHTTP(w, r)
		})
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", ts.cfg.Port),
		Handler:      logger(ts.Router),
		ReadTimeout:  time.Duration(ts.cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(ts.cfg.WriteTimeout) * time.Second,
	}
	ts.Server = srv

	ts.logger.For(ctx).Info(fmt.Sprintf("starting Telemetry server on port %d", ts.cfg.Port))
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		ts.logger.For(ctx).With(zap.Error(err)).Error("failed to start Telemetry server")
		return err
	}
	return nil
}

func (ts *TelemetryServer) Stop(ctx context.Context) error {
	ts.logger.For(ctx).Info("stopping Telemetry server")
	err := ts.Server.Close()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		ts.logger.For(ctx).With(zap.Error(err)).Error("failed to stop Telemetry server")
		return err
	}
	ts.logger.For(ctx).Info("Telemetry server stopped")
	return nil
}
