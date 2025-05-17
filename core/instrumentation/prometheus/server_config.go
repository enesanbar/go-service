package prometheus

import (
	"fmt"

	"github.com/enesanbar/go-service/core/config"
)

const (
	Port        = "port"
	PortDefault = 9092

	ReadTimeout        = "read_timeout"
	ReadTimeoutDefault = 10

	WriteTimeout        = "write_timeout"
	WriteTimeoutDefault = 20

	GracefulStopTimeoutSeconds = "gracefulStopTimeoutSeconds"
	GracefulStopTimeoutDefault = 10
)

type TelemetryServerConfig struct {
	Port                       int
	ReadTimeout                int
	WriteTimeout               int
	GracefulStopTimeoutSeconds int
}

func NewTelemetryServerConfig(cfg config.Config) *TelemetryServerConfig {
	key := "server.telemetry.%s"

	port := cfg.GetInt(fmt.Sprintf(key, Port))
	if port == 0 {
		port = PortDefault
	}

	readTimeout := cfg.GetInt(fmt.Sprintf(key, ReadTimeout))
	if readTimeout == 0 {
		readTimeout = ReadTimeoutDefault
	}

	writeTimeout := cfg.GetInt(fmt.Sprintf(key, WriteTimeout))
	if writeTimeout == 0 {
		writeTimeout = WriteTimeoutDefault
	}

	gracefulStopTimeout := cfg.GetInt(fmt.Sprintf(key, GracefulStopTimeoutSeconds))
	if gracefulStopTimeout == 0 {
		gracefulStopTimeout = GracefulStopTimeoutDefault
	}

	return &TelemetryServerConfig{
		Port:                       port,
		ReadTimeout:                readTimeout,
		WriteTimeout:               writeTimeout,
		GracefulStopTimeoutSeconds: gracefulStopTimeout,
	}
}
