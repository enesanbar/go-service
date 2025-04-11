package router

import (
	"net/http"

	"github.com/enesanbar/go-service/healthchecker"
	"github.com/enesanbar/go-service/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HealthCheckHandler struct {
	healthChecker *healthchecker.HealthChecker
	logger        log.Factory
}

func NewHealthCheckHandler(
	healthchecker *healthchecker.HealthChecker,
	logger log.Factory,
) *HealthCheckHandler {
	return &HealthCheckHandler{healthChecker: healthchecker, logger: logger}
}

func (h *HealthCheckHandler) Handle(c echo.Context) error {
	r := h.healthChecker.Run(c.Request().Context())
	if !r.Success {
		h.logger.For(c.Request().Context()).Error("health check failed", zap.Any("message", r.ProbesResults))
		return c.JSON(http.StatusInternalServerError, r)
	}

	return c.JSON(http.StatusOK, r)
}
