package prometheus

import (
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

func (ts *TelemetryServer) Start() error {
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

	ts.logger.Bg().Info(fmt.Sprintf("starting Telemetry server on port %d", ts.cfg.Port))
	return srv.ListenAndServe()
}

func (ts *TelemetryServer) Stop() error {
	ts.logger.Bg().Info("stopping Telemetry server")
	return ts.Server.Close()
}
