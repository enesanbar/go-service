package producer

import "context"

type Producer interface {
	Publish(ctx context.Context, messageName string, message any) error
}
