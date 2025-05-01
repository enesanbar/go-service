package grpc

import (
	"fmt"

	"github.com/enesanbar/go-service/config"
)

const (
	Port        = "port"
	PortDefault = 50051

	GracefulStopTimeoutSeconds = "graceful_stop_timeout_seconds"
	GracefulStopTimeoutDefault = 10
)

type ServerConfig struct {
	Port                       int
	GracefulStopTimeoutSeconds int
}

func NewServerConfig(cfg config.Config) *ServerConfig {
	key := "server.grpc.%s"

	property := fmt.Sprintf(key, Port)
	port := cfg.GetInt(property)
	if port == 0 {
		port = PortDefault
	}

	property = fmt.Sprintf(key, GracefulStopTimeoutSeconds)
	gracefulStopTimeout := cfg.GetInt(property)
	if gracefulStopTimeout == 0 {
		gracefulStopTimeout = GracefulStopTimeoutDefault
	}

	return &ServerConfig{
		Port:                       port,
		GracefulStopTimeoutSeconds: gracefulStopTimeout,
	}
}
