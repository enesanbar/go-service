package grpc

import (
	"fmt"

	"github.com/enesanbar/go-service/config"
)

const (
	Port        = "port"
	PortDefault = 50051

	GracefulStopTimeoutSeconds = "gracefulStopTimeoutSeconds"
	GracefulStopTimeoutDefault = 10

	KeepAliveMinTimeSeconds        = "keepalive.minTimeSeconds"
	KeepAliveMinTimeSecondsDefault = 300 // 5 minutes

	KeepAlivePermitWithoutStream        = "keepalive.permitWithoutStream"
	KeepAlivePermitWithoutStreamDefault = false

	KeepAliveMaxConnectionIdleSeconds        = "keepalive.maxConnectionIdleSeconds"
	KeepAliveMaxConnectionIdleSecondsDefault = 15

	KeepAliveMaxConnectionAgeSeconds        = "keepalive.maxConnectionAgeSeconds"
	KeepAliveMaxConnectionAgeSecondsDefault = 30

	KeepAliveMaxConnectionAgeGraceSeconds        = "keepalive.maxConnectionAgeGraceSeconds"
	KeepAliveMaxConnectionAgeGraceSecondsDefault = 5

	KeepAliveTimeSeconds        = "keepalive.timeSeconds"
	KeepAliveTimeSecondsDefault = 7200 // 2 hours

	KeepAliveTimeoutSeconds        = "keepalive.timeoutSeconds"
	KeepAliveTimeoutSecondsDefault = 20 // 20 seconds
)

type KeepAlive struct {
	MinTimeSeconds               int  `json:"minTimeSeconds" yaml:"minTimeSeconds"`                             // If a client pings more than once every specified seconds, terminate the connection
	PermitWithoutStream          bool `json:"permitWithoutStream" yaml:"permitWithoutStream"`                   // Allow pings even when there are no active streams
	MaxConnectionIdleSeconds     int  `json:"maxConnectionIdleSeconds" yaml:"maxConnectionIdleSeconds"`         // If a client is idle for specified seconds, send a GOAWAY
	MaxConnectionAgeSeconds      int  `json:"maxConnectionAgeSeconds" yaml:"maxConnectionAgeSeconds"`           // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGraceSeconds int  `json:"maxConnectionAgeGraceSeconds" yaml:"maxConnectionAgeGraceSeconds"` // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	TimeSeconds                  int  `json:"timeSeconds" yaml:"timeSeconds"`                                   // Ping the client if it is idle for specified seconds to ensure the connection is still active
	TimeoutSeconds               int  `json:"timeoutSeconds" yaml:"timeoutSeconds"`                             // // Wait 1 second for the ping ack before assuming the connection is dead
}

type ServerConfig struct {
	Port                       int       `json:"port" yaml:"port"`
	GracefulStopTimeoutSeconds int       `json:"gracefulStopTimeoutSeconds" yaml:"gracefulStopTimeoutSeconds"`
	KeepAlive                  KeepAlive `json:"keepalive" yaml:"keepalive"`
}

func NewServerConfig(cfg config.Config) *ServerConfig {
	key := "server.grpc.%s"

	// note: doing it manually. because viper only take environment variables into account
	// when the config is read individually with cfg.Get*()
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

	property = fmt.Sprintf(key, KeepAliveMinTimeSeconds)
	keepAliveMinTime := cfg.GetInt(property)
	if keepAliveMinTime == 0 {
		keepAliveMinTime = KeepAliveMinTimeSecondsDefault
	}

	property = fmt.Sprintf(key, KeepAlivePermitWithoutStream)
	keepAlivePermitWithoutStream := cfg.GetBool(property)

	property = fmt.Sprintf(key, KeepAliveMaxConnectionIdleSeconds)
	keepAliveMaxConnectionIdle := cfg.GetInt(property)

	property = fmt.Sprintf(key, KeepAliveMaxConnectionAgeSeconds)
	keepAliveMaxConnectionAge := cfg.GetInt(property)

	property = fmt.Sprintf(key, KeepAliveMaxConnectionAgeGraceSeconds)
	keepAliveMaxConnectionAgeGrace := cfg.GetInt(property)

	property = fmt.Sprintf(key, KeepAliveTimeSeconds)
	keepAliveTime := cfg.GetInt(property)
	if keepAliveTime == 0 {
		keepAliveTime = KeepAliveTimeSecondsDefault
	}

	property = fmt.Sprintf(key, KeepAliveTimeoutSeconds)
	keepAliveTimeout := cfg.GetInt(property)
	if keepAliveTimeout == 0 {
		keepAliveTimeout = KeepAliveTimeoutSecondsDefault
	}

	keepAlive := KeepAlive{
		MinTimeSeconds:      keepAliveMinTime,
		PermitWithoutStream: keepAlivePermitWithoutStream,
		TimeSeconds:         keepAliveTime,
		TimeoutSeconds:      keepAliveTimeout,
	}

	// the default values for MaxConnectionIdle, MaxConnectionAge, and MaxConnectionAgeGrace are infinity
	// so we only set them if they are not 0
	if keepAliveMaxConnectionIdle != 0 {
		keepAlive.MaxConnectionIdleSeconds = keepAliveMaxConnectionIdle
	}
	if keepAliveMaxConnectionAge != 0 {
		keepAlive.MaxConnectionAgeSeconds = keepAliveMaxConnectionAge
	}
	if keepAliveMaxConnectionAgeGrace != 0 {
		keepAlive.MaxConnectionAgeGraceSeconds = keepAliveMaxConnectionAgeGrace
	}

	return &ServerConfig{
		Port:                       port,
		GracefulStopTimeoutSeconds: gracefulStopTimeout,
		KeepAlive:                  keepAlive,
	}
}
