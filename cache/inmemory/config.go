package inmemory

import (
	"fmt"
	"time"

	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
)

const (
	PropertyExpiration      = "default-expiration"
	PropertyCleanupInterval = "default-expiration"
)

type Config struct {
	Expiration      time.Duration
	CleanupInterval time.Duration
}

func New(cfg config.Config, logger log.Factory, prefix string) (*Config, error) {
	keyTemplate := "%s.%s"
	property := fmt.Sprintf(keyTemplate, prefix, PropertyExpiration)

	expiration := cfg.GetInt(property)
	if expiration == 0 {
		return nil, config.NewMissingPropertyError(property)
	}

	property = fmt.Sprintf(keyTemplate, prefix, PropertyCleanupInterval)

	cleanupInterval := cfg.GetInt(property)
	if cleanupInterval == 0 {
		panic(fmt.Sprintf("please set '%s.%s' in consul", prefix, property))
	}

	return &Config{
		Expiration:      time.Duration(expiration) * time.Minute,
		CleanupInterval: time.Duration(cleanupInterval) * time.Minute,
	}, nil
}
