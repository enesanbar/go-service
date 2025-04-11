package rabbitmq

import (
	"fmt"

	"github.com/enesanbar/go-service/config"
)

const (
	PropertyHost     = "host"
	PropertyPort     = "port"
	PropertyUsername = "username"
	PropertyPassword = "password"
)

type Config struct {
	Name string
	Port string
	Host string
	User string
	Pass string
}

func NewConfig(cfg config.Config, name string) (*Config, error) {
	keyTemplate := "datasources.rabbitmq.%s.%s"

	property := fmt.Sprintf(keyTemplate, name, PropertyHost)
	host := cfg.GetString(property)
	if host == "" {
		host = "localhost" // Default RabbitMQ host
	}

	port := cfg.GetString(fmt.Sprintf(keyTemplate, name, "port"))
	if port == "" {
		port = "5672" // Default RabbitMQ port
	}

	property = fmt.Sprintf(keyTemplate, name, PropertyUsername)
	username := cfg.GetString(property)
	if username == "" {
		return nil, config.NewMissingPropertyError(property)
	}

	password := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyPassword))
	if password == "" {
		return nil, config.NewMissingPropertyError(property)
	}

	return &Config{
		Host: host,
		Port: port,
		User: username,
		Pass: password,
	}, nil
}
