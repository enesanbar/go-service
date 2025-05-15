package inmemory

import (
	"fmt"
	"time"

	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
)

const (
	ExpirationKey     = "default-expiration"
	ExpirationDefault = 5

	CleanupIntervalKey     = "default-expiration"
	CleanupIntervalDefault = 10
)

type Config struct {
	Expiration      time.Duration
	CleanupInterval time.Duration
}

func NewConfig(cfg config.Config, logger log.Factory) (*Config, error) {
	prefix := "cache.inmemory"
	keyTemplate := "%s.%s"
	property := fmt.Sprintf(keyTemplate, prefix, ExpirationKey)

	expiration := cfg.GetInt(property)
	if expiration == 0 {
		expiration = ExpirationDefault
	}

	property = fmt.Sprintf(keyTemplate, prefix, CleanupIntervalKey)

	cleanupInterval := cfg.GetInt(property)
	if cleanupInterval == 0 {
		cleanupInterval = CleanupIntervalDefault
	}

	return &Config{
		Expiration:      time.Duration(expiration) * time.Minute,
		CleanupInterval: time.Duration(cleanupInterval) * time.Minute,
	}, nil
}
