package grpc

import (
	"encoding/json"
	"fmt"

	"github.com/enesanbar/go-service/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type RetryPolicy struct {
	MaxAttempts          int      `json:"maxAttempts"`
	InitialBackoff       string   `json:"initialBackoff"`
	MaxBackoff           string   `json:"maxBackoff"`
	BackoffMultiplier    float64  `json:"backoffMultiplier"`
	RetryableStatusCodes []string `json:"retryableStatusCodes"`
}

type MethodConfig struct {
	Name        []map[string]string `json:"name"`
	RetryPolicy RetryPolicy         `json:"retryPolicy"`
	Timeout     string              `json:"timeout"`
}

type ServiceConfig struct {
	MethodConfig []MethodConfig `json:"methodConfig"`
}

// Params to inject Viper-based config
type GRPCClientOptionFactoryParams struct {
	fx.In

	Config config.Config
}

// Provide this in fx.Module
func NewGRPCClientOptionFactory(Config config.Config) func(serviceName string) grpc.DialOption {
	return func(serviceName string) grpc.DialOption {
		configKey := fmt.Sprintf("client.grpc.%s.serviceConfig", serviceName)

		var serviceConfig ServiceConfig
		if Config.IsSet(configKey) {
			err := Config.UnmarshalKey(configKey, &serviceConfig)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal retry config for %s: %w", serviceName, err))
			}
		} else {
			serviceConfig = ServiceConfig{
				MethodConfig: []MethodConfig{{
					Name:    []map[string]string{{"service": ""}},
					Timeout: "10s",
					RetryPolicy: RetryPolicy{
						MaxAttempts:       3,
						InitialBackoff:    "0.1s",
						MaxBackoff:        "5s",
						BackoffMultiplier: 1.5,
						RetryableStatusCodes: []string{
							"UNAVAILABLE",
							"DEADLINE_EXCEEDED",
						},
					},
				}},
			}
		}

		configBytes, err := json.Marshal(serviceConfig)
		if err != nil {
			panic(fmt.Errorf("failed to marshal service config: %w", err))
		}

		return grpc.WithDefaultServiceConfig(string(configBytes))
	}
}
