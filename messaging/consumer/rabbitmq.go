package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/messaging/messages"
	"github.com/enesanbar/go-service/messaging/rabbitmq"
	"github.com/enesanbar/go-service/wiring"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RabbitMQConsumerParams struct {
	fx.In

	Logger          log.Factory
	Queues          map[string]*rabbitmq.Queue
	Channels        map[string]*rabbitmq.Channel
	MessageHandlers map[string]MessageHandler
}

type RabbitMQQueueConsumer struct {
	logger          log.Factory
	Channel         *rabbitmq.Channel
	Queue           *rabbitmq.Queue
	Channels        map[string]*rabbitmq.Channel
	Queues          map[string]*rabbitmq.Queue
	MessageHandlers map[string]MessageHandler
}

// NewRabbitMQConsumer creates a pointer to the new instance of the RabbitMQQueueConsumer
// and a runnable group that can be used to start and stop the consumer
func NewRabbitMQConsumer(p RabbitMQConsumerParams) (wiring.RunnableGroup, *RabbitMQQueueConsumer) {
	consumer := &RabbitMQQueueConsumer{
		logger:          p.Logger,
		MessageHandlers: p.MessageHandlers,
		Channels:        p.Channels,
		Queues:          p.Queues,
	}

	// set the default channel and queue if it exists in Channels and Queues
	if len(p.Channels) > 0 {
		for name, channel := range p.Channels {
			if name == "default" {
				consumer.Channel = channel
				break
			}
		}
	}

	if len(p.Queues) > 0 {
		for name, queue := range p.Queues {
			if name == "default" {
				consumer.Queue = queue
				break
			}
		}
	}

	return wiring.RunnableGroup{
		Runnable: consumer,
	}, consumer
}

func (h *RabbitMQQueueConsumer) Start() error {
	if h.Channel == nil {
		return fmt.Errorf("channel is not set, check your configuration")
	}
	if h.Queue == nil {
		return fmt.Errorf("queue is not set, check your configuration")
	}

	// add recovery logic to the channel when channel/connection is closed
	msgs, err := h.Channel.Channel.Consume(
		h.Queue.Queue.Name, // queue
		"",                 // consumer
		true,               // auto ack
		false,              // exclusive
		false,              // no local
		false,              // no wait
		nil,                // args
	)
	if err != nil {
		return fmt.Errorf("error starting RabbitMQ consumer (%w)", err)
	}

	go func() {
		for d := range msgs {
			// TODO: throttle the message processing with worker pool
			go func(d amqp091.Delivery) {
				message := messages.Message[any]{}
				err := json.Unmarshal(d.Body, &message)
				if err != nil {
					h.logger.Bg().With(zap.Error(err)).Error("Failed to unmarshal message")
					return
				}

				key := fmt.Sprintf("%s-%s", h.Queue.Config.Name, message.Metadata.MessageName)
				handler, ok := h.MessageHandlers[key]
				if !ok {
					h.logger.Bg().With(zap.String("messageName", message.Metadata.MessageName)).Error("no handler found for message")
					return
				}

				err = handler.Handle(message)
				if err != nil {
					h.logger.Bg().With(zap.Error(err)).Error("failed to handle message")
				}
			}(d)
		}
		h.logger.Bg().Info("RabbitMQ consumer stopped")
	}()
	h.logger.Bg().Info(fmt.Sprintf("RabbitMQ consumer started for queue %s", h.Queue.Config.Name))
	return nil
}

func (h *RabbitMQQueueConsumer) Stop() error {
	h.logger.Bg().Info("RabbitMQ consumer stopped")
	return nil
}

func (h *RabbitMQQueueConsumer) SetChannel(channelName string) {
	channel, ok := h.Channels[channelName]
	if !ok {
		h.logger.Bg().Error(fmt.Sprintf("channel %s not found, check your configuration", channelName))
	}
	h.Channel = channel
}
func (h *RabbitMQQueueConsumer) SetQueue(queueName string) {
	queue, ok := h.Queues[queueName]
	if !ok {
		h.logger.Bg().Error(fmt.Sprintf("queue %s not found, check your configuration", queueName))
	}
	h.Queue = queue
}
