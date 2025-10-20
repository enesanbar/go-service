package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/enesanbar/go-service/core/healthchecker"
	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/log"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthCheckHandler struct {
	healthChecker     *healthchecker.HealthChecker
	logger            log.Factory
	GRPCServer        *Server
	HealthCheckServer *health.Server
}

func NewHealthCheckHandler(
	healthchecker *healthchecker.HealthChecker,
	logger log.Factory,
	grpcServer *Server,
) *HealthCheckHandler {
	healthcheck := health.NewServer()
	healthcheck.SetServingStatus(info.ServiceName, healthpb.HealthCheckResponse_SERVING)

	logger.Bg().Info("health check server registered in grpc server")
	healthgrpc.RegisterHealthServer(grpcServer.Server, healthcheck)

	hc := &HealthCheckHandler{
		logger:            logger,
		healthChecker:     healthchecker,
		GRPCServer:        grpcServer,
		HealthCheckServer: healthcheck,
	}
	// TODO: for now, let's invoke here. it should be invoked from the main app lifecycle as runnable
	hc.Handle(context.Background())
	return hc
}

func (h *HealthCheckHandler) Handle(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r := h.healthChecker.Run(ctx)
				if !r.Success {
					h.HealthCheckServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
					h.logger.For(ctx).Error("health check failed", zap.Any("message", r.ProbesResults))
					continue
				}

				// Health check also supports checking serving status for a specific service
				// in this case, we pass the service name as empty and setting the status for empty service name
				h.HealthCheckServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
			}
		}
	}()
}
