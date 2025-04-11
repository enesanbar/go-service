package http

import (
	"fmt"

	"github.com/enesanbar/go-service/config"
)

const (
	Port        = "port"
	PortDefault = 9090

	ReadTimeout        = "read_timeout"
	ReadTimeoutDefault = 10

	WriteTimeout        = "write_timeout"
	WriteTimeoutDefault = 20
)

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

func NewServerConfig(cfg config.Config) *ServerConfig {
	key := "server.%s"

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

	return &ServerConfig{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
