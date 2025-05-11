package inmemory

import (
	"context"
	"errors"

	"github.com/patrickmn/go-cache"

	"github.com/enesanbar/go-service/cache/metrics"
	"github.com/enesanbar/go-service/core/log"
)

var ErrKeyNotFound = errors.New("'key' not found in cache")

// Cache is an in memory cache implementation of Cache interface
type Cache struct {
	cfg          *Config
	cache        *cache.Cache
	log          log.Factory
	instrumentor *metrics.Instrumentor
}

// NewInMemoryCache returns a pointer to the instance of Cache
func NewInMemoryCache(cfg *Config, log log.Factory, instrumentor *metrics.Instrumentor) *Cache {
	return &Cache{
		cfg:          cfg,
		cache:        cache.New(cfg.Expiration, cfg.CleanupInterval),
		log:          log,
		instrumentor: instrumentor,
	}
}

// Set sets a value in cache
func (c *Cache) Set(_ context.Context, key string, value interface{}) {
	c.cache.Set(key, value, c.cfg.Expiration)
}

// Get gets a value from the cache
func (c *Cache) Get(_ context.Context, key string) (interface{}, error) {
	cached, found := c.cache.Get(key)
	if found {
		c.instrumentor.Hit(key, "inmemory")
		return cached, nil
	}

	c.instrumentor.Miss(key, "inmemory")
	return nil, ErrKeyNotFound
}
