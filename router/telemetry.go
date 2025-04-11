package router

import (
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/wiring"

	"github.com/enesanbar/go-service/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TelemetryServer struct {
	Router     *http.ServeMux
	logger     log.Factory
	Server     io.Closer
	BaseConfig *config.Base
}

func NewTelemetryServer(logger log.Factory, baseConfig *config.Base) (wiring.RunnableGroup, *TelemetryServer) {
	telemetryRouter := http.NewServeMux()

	telemetryRouter.Handle("/metrics", promhttp.Handler())

	server := &TelemetryServer{
		Router:     telemetryRouter,
		logger:     logger,
		BaseConfig: baseConfig,
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
		Addr:         ":9092",
		Handler:      logger(ts.Router),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	ts.Server = srv

	ts.logger.Bg().Info("starting Telemetry server on 9092...")
	return srv.ListenAndServe()
}

func (ts *TelemetryServer) Stop() error {
	ts.logger.Bg().Info("stopping Telemetry server")
	return ts.Server.Close()
}
