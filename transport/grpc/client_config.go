package grpc

// import (
// 	"fmt"

// 	"github.com/enesanbar/go-service/config"
// )

// const (
// 	KeepAlivePermitWithoutStream        = "keepalive.permitWithoutStream"
// 	KeepAlivePermitWithoutStreamDefault = false

// 	KeepAliveTimeSeconds        = "keepalive.timeSeconds"
// 	KeepAliveTimeSecondsDefault = 7200 // 2 hours

// 	KeepAliveTimeoutSeconds        = "keepalive.timeoutSeconds"
// 	KeepAliveTimeoutSecondsDefault = 20 // 20 seconds
// )

// type KeepAliveClient struct {
// 	MinTimeSeconds               int  `json:"minTimeSeconds" yaml:"minTimeSeconds"`                             // If a client pings more than once every specified seconds, terminate the connection
// 	PermitWithoutStream          bool `json:"permitWithoutStream" yaml:"permitWithoutStream"`                   // Allow pings even when there are no active streams
// 	MaxConnectionIdleSeconds     int  `json:"maxConnectionIdleSeconds" yaml:"maxConnectionIdleSeconds"`         // If a client is idle for specified seconds, send a GOAWAY
// 	MaxConnectionAgeSeconds      int  `json:"maxConnectionAgeSeconds" yaml:"maxConnectionAgeSeconds"`           // If any connection is alive for more than 30 seconds, send a GOAWAY
// 	MaxConnectionAgeGraceSeconds int  `json:"maxConnectionAgeGraceSeconds" yaml:"maxConnectionAgeGraceSeconds"` // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
// 	TimeSeconds                  int  `json:"timeSeconds" yaml:"timeSeconds"`                                   // Ping the client if it is idle for specified seconds to ensure the connection is still active
// 	TimeoutSeconds               int  `json:"timeoutSeconds" yaml:"timeoutSeconds"`                             // // Wait 1 second for the ping ack before assuming the connection is dead
// }

// type ClientConfig struct {
// 	Port                       int             `json:"port" yaml:"port"`
// 	GracefulStopTimeoutSeconds int             `json:"gracefulStopTimeoutSeconds" yaml:"gracefulStopTimeoutSeconds"`
// 	KeepAlive                  KeepAliveClient `json:"keepalive" yaml:"keepalive"`
// 	TLS                        bool            `json:"tls" yaml:"tls"` // if true, use TLS
// }

// func NewClientConfig(cfg config.Config) *ClientConfig {
// 	key := "server.grpc.%s"

// 	property := fmt.Sprintf(key, KeepAlivePermitWithoutStream)
// 	keepAlivePermitWithoutStream := cfg.GetBool(property)

// 	property = fmt.Sprintf(key, KeepAliveTimeSeconds)
// 	keepAliveTime := cfg.GetInt(property)
// 	if keepAliveTime == 0 {
// 		keepAliveTime = KeepAliveTimeSecondsDefault
// 	}

// 	property = fmt.Sprintf(key, KeepAliveTimeoutSeconds)
// 	keepAliveTimeout := cfg.GetInt(property)
// 	if keepAliveTimeout == 0 {
// 		keepAliveTimeout = KeepAliveTimeoutSecondsDefault
// 	}

// 	keepAlive := KeepAliveClient{
// 		PermitWithoutStream: keepAlivePermitWithoutStream,
// 		TimeSeconds:         keepAliveTime,
// 		TimeoutSeconds:      keepAliveTimeout,
// 	}

// 	return &ClientConfig{
// 		KeepAlive: keepAlive,
// 	}
// }
