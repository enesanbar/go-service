package consumer

import (
	"github.com/enesanbar/go-service/messaging/messages"
)

type MessageHandler interface {
	Handle(message messages.Message[any]) error
	Properties() MessageProperties
}
