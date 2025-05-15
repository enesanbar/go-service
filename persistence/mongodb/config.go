package mongodb

import (
	"fmt"
	"time"

	"github.com/enesanbar/go-service/core/config"
)

const (
	PropertyHost           = "host"
	PropertyName           = "name"
	PropertyAuthDB         = "auth_db"
	PropertyReplicaSetName = "replica_set"
)

type Config struct {
	Host                  string
	AuthDB                string
	User                  string
	Pass                  string
	Name                  string
	ReplicaSetName        string
	Timeout               time.Duration
	MinPoolSize           uint64
	MaxPoolSize           uint64
	MaxConnectionIdletime time.Duration
}

func NewConfig(cfg config.Config, prefix string) (*Config, error) {
	keyTemplate := "%s.%s"

	property := fmt.Sprintf(keyTemplate, prefix, PropertyHost)
	host := cfg.GetString(property)
	if host == "" {
		return nil, config.NewMissingPropertyError(property)
	}

	property = fmt.Sprintf(keyTemplate, prefix, PropertyName)
	name := cfg.GetString(property)
	if name == "" {
		return nil, config.NewMissingPropertyError(property)
	}

	property = fmt.Sprintf(keyTemplate, prefix, PropertyReplicaSetName)
	rsName := cfg.GetString(property)

	property = fmt.Sprintf(keyTemplate, prefix, PropertyAuthDB)
	authDB := cfg.GetString(property)

	// TODO: Check rest of the config and determine default values
	return &Config{
		Host:                  host,
		Name:                  name,
		ReplicaSetName:        rsName,
		Timeout:               time.Duration(cfg.GetInt("datasources.mongo.default.timeout")) * time.Second,
		AuthDB:                authDB,
		User:                  cfg.GetString("datasources.mongo.default.username"),
		Pass:                  cfg.GetString("datasources.mongo.default.password"),
		MaxConnectionIdletime: time.Duration(cfg.GetInt("datasources.mongo.default.max-conn-idle-time")) * time.Second,
		MaxPoolSize:           uint64(cfg.GetInt("datasources.mongo.default.max-pool-size")),
		MinPoolSize:           uint64(cfg.GetInt("datasources.mongo.default.min-pool-size")),
	}, nil
}
