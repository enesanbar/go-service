package grpc

import (
	"context"

	"github.com/enesanbar/go-service/healthchecker"
	"github.com/enesanbar/go-service/info"
	"github.com/enesanbar/go-service/log"
	"go.uber.org/zap"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthCheckHandler struct {
	healthChecker     *healthchecker.HealthChecker
	logger            log.Factory
	GRPCServer        *GRPCServer
	HealthCheckServer *health.Server
}

func NewHealthCheckHandler(
	healthchecker *healthchecker.HealthChecker,
	logger log.Factory,
	grpcServer *GRPCServer,
) *HealthCheckHandler {
	healthcheck := health.NewServer()
	healthcheck.SetServingStatus(info.ServiceName, healthpb.HealthCheckResponse_SERVING)

	logger.Bg().Info("health check server registered in grpc server")
	healthgrpc.RegisterHealthServer(grpcServer.Server, healthcheck)

	return &HealthCheckHandler{
		logger:            logger,
		healthChecker:     healthchecker,
		GRPCServer:        grpcServer,
		HealthCheckServer: healthcheck,
	}
}

func (h *HealthCheckHandler) Handle(ctx context.Context) {
	r := h.healthChecker.Run(ctx)
	if !r.Success {
		h.logger.For(ctx).Error("health check failed", zap.Any("message", r.ProbesResults))
		h.HealthCheckServer.SetServingStatus(info.ServiceName, healthpb.HealthCheckResponse_NOT_SERVING)
		return
	}

	h.HealthCheckServer.SetServingStatus(info.ServiceName, healthpb.HealthCheckResponse_SERVING)
}
