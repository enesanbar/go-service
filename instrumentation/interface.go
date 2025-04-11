package instrumentation

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// HTTP application
type HTTP struct {
	Handler      string
	Host         string
	Method       string
	StatusCode   string
	StartedAt    time.Time
	Duration     float64
	RequestSize  int64
	ResponseSize int64
}

// NewHTTP create a new HTTP app
func NewHTTP(handler string, method string) *HTTP {
	return &HTTP{
		Handler:   handler,
		Method:    method,
		StartedAt: time.Now(),
	}
}

// Finished app finished
func (h *HTTP) Finished(c echo.Context) {
	h.Duration = time.Since(h.StartedAt).Seconds()
	h.StatusCode = strconv.Itoa(c.Response().Status)
	h.Host = c.Request().Host
	h.RequestSize = c.Request().ContentLength
	h.ResponseSize = c.Response().Size
}
