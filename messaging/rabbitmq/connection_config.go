package rabbitmq

import (
	"fmt"

	"github.com/enesanbar/go-service/core/config"
)

const (
	PropertyHost     = "host"
	PropertyPort     = "port"
	PropertyUsername = "username"
	PropertyPassword = "password"
)

type ConnectionConfig struct {
	Name string
	Port string
	Host string
	User string
	Pass string
}

func NewConnectionConfig(cfg config.Config, name string) (*ConnectionConfig, error) {
	keyTemplate := "rabbitmq.connections.%s.%s"

	host := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyHost))
	if host == "" {
		host = "localhost" // Default RabbitMQ host
	}

	port := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyPort))
	if port == "" {
		port = "5672" // Default RabbitMQ port
	}

	property := fmt.Sprintf(keyTemplate, name, PropertyUsername)
	username := cfg.GetString(property)
	if username == "" {
		return nil, config.NewMissingPropertyError(property)
	}

	property = fmt.Sprintf(keyTemplate, name, PropertyPassword)
	password := cfg.GetString(property)
	if password == "" {
		return nil, config.NewMissingPropertyError(property)
	}

	return &ConnectionConfig{
		Name: name,
		Host: host,
		Port: port,
		User: username,
		Pass: password,
	}, nil
}
