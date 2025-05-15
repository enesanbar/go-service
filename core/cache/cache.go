package cache

import (
	"context"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	Invalidate(ctx context.Context, key string) error
}
