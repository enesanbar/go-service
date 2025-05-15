package rest

import (
	"fmt"

	"github.com/enesanbar/go-service/core/config"
)

const (
	Port        = "port"
	PortDefault = 9090

	ReadTimeout        = "read_timeout"
	ReadTimeoutDefault = 10

	WriteTimeout        = "write_timeout"
	WriteTimeoutDefault = 20

	GracefulStopTimeoutSeconds = "graceful_stop_timeout_seconds"
	GracefulStopTimeoutDefault = 10
)

type ServerConfig struct {
	Port                       int
	ReadTimeout                int
	WriteTimeout               int
	GracefulStopTimeoutSeconds int
}

func NewConfig(cfg config.Config) *ServerConfig {
	key := "server.http.%s"

	property := fmt.Sprintf(key, Port)
	port := cfg.GetInt(property)
	if port == 0 {
		port = PortDefault
	}

	property = fmt.Sprintf(key, ReadTimeout)
	readTimeout := cfg.GetInt(property)
	if port == 0 {
		readTimeout = ReadTimeoutDefault
	}

	property = fmt.Sprintf(key, WriteTimeout)
	writeTimeout := cfg.GetInt(property)
	if port == 0 {
		writeTimeout = WriteTimeoutDefault
	}

	gracefulStopTimeout := cfg.GetInt(property)
	if gracefulStopTimeout == 0 {
		gracefulStopTimeout = GracefulStopTimeoutDefault
	}

	return &ServerConfig{
		Port:                       port,
		ReadTimeout:                readTimeout,
		WriteTimeout:               writeTimeout,
		GracefulStopTimeoutSeconds: gracefulStopTimeout,
	}
}
