package wiring

import "context"

type Connection interface {
	Name() string
	Close(ctx context.Context) error
	Start(ctx context.Context) error
}
