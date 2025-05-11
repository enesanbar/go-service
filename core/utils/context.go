package utils

import (
	"context"
)

var (
	ContextKeyRequestID = NewContextKey("X-Request-Id")
	ContextKeyUsername  = NewContextKey("username")
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func NewContextKey(key string) ContextKey {
	return ContextKey(key)
}

// GetValueFromContext gets the value from the context.
func GetValueFromContext(ctx context.Context, key ContextKey) (string, bool) {
	caller, ok := ctx.Value(key).(string)
	return caller, ok
}
