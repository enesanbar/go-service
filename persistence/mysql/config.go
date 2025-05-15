package mysql

import (
	"fmt"

	"github.com/enesanbar/go-service/core/config"
)

const (
	PropertyName                  = "name"
	PropertyDatabase              = "database"
	PropertyHost                  = "host"
	PropertyPort                  = "port"
	PropertyUsername              = "username"
	PropertyPassword              = "password"
	PropertyTimeout               = "timeout"
	PropertyMaxIdleConnections    = "maxIdleConnections"
	PropertyMaxOpenConnections    = "maxOpenConnections"
	PropertyMaxConnectionLifetime = "maxConnectionLifetime"
	PropertyMaxConnectionIdleTime = "maxConnectionIdletime"

	// Default values
	DefaultHost                  = "localhost"
	DefaultPort                  = 3306
	DefaultTimeout               = 30
	DefaultMaxIdleConnections    = 10
	DefaultMaxOpenConnections    = 100
	DefaultMaxConnectionLifetime = 3600
	DefaultMaxConnectionIdleTime = 300
)

type Config struct {
	Name                  string
	Database              string
	Host                  string
	Port                  int
	User                  string
	Pass                  string
	Timeout               int
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifetime int
	MaxConnectionIdleTime int
}

func NewConfig(cfg config.Config, name string) (*Config, error) {
	keyTemplate := "datasources.mysql.%s.%s"

	host := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyHost))
	if host == "" {
		host = DefaultHost
	}

	port := cfg.GetInt(fmt.Sprintf(keyTemplate, name, "port"))
	if port == 0 {
		port = DefaultPort
	}

	database := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyDatabase))
	username := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyUsername))
	password := cfg.GetString(fmt.Sprintf(keyTemplate, name, PropertyPassword))

	timeout := cfg.GetInt(fmt.Sprintf(keyTemplate, name, PropertyTimeout))
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	maxIdleConnections := cfg.GetInt(fmt.Sprintf(keyTemplate, name, PropertyMaxIdleConnections))
	if maxIdleConnections == 0 {
		maxIdleConnections = DefaultMaxIdleConnections
	}
	maxOpenConnections := cfg.GetInt(fmt.Sprintf(keyTemplate, name, PropertyMaxOpenConnections))
	if maxOpenConnections == 0 {
		maxOpenConnections = DefaultMaxOpenConnections
	}
	maxConnectionLifetime := cfg.GetInt(fmt.Sprintf(keyTemplate, name, PropertyMaxConnectionLifetime))
	if maxConnectionLifetime == 0 {
		maxConnectionLifetime = DefaultMaxConnectionLifetime
	}
	maxConnectionIdleTime := cfg.GetInt(fmt.Sprintf(keyTemplate, name, PropertyMaxConnectionIdleTime))
	if maxConnectionIdleTime == 0 {
		maxConnectionIdleTime = DefaultMaxConnectionIdleTime
	}

	return &Config{
		Name:                  name,
		Database:              database,
		Host:                  host,
		Port:                  port,
		User:                  username,
		Pass:                  password,
		Timeout:               timeout,
		MaxIdleConnections:    maxIdleConnections,
		MaxOpenConnections:    maxOpenConnections,
		MaxConnectionLifetime: maxConnectionLifetime,
		MaxConnectionIdleTime: maxConnectionIdleTime,
	}, nil
}
