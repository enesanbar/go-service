package cache

import (
	"context"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{})
	Get(ctx context.Context, key string) (interface{}, error)
}
