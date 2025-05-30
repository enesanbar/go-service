package consumer

import (
	"context"

	"github.com/enesanbar/go-service/core/messaging/messages"
)

type MessageHandler interface {
	Handle(ctx context.Context, message messages.Message[any]) error
	Properties() MessageProperties
	GetMessageType() any
}
