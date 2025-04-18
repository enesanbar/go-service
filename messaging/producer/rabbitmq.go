package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/enesanbar/go-service/info"
	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/messaging/rabbitmq"
	"github.com/rabbitmq/amqp091-go"

	"go.uber.org/fx"
)

// TODO: Move this to a common package
type RabbitMQProducer struct {
	Logger   log.Factory
	Channel  *rabbitmq.Channel
	Channels map[string]*rabbitmq.Channel
}

type RabbitMQProducerParams struct {
	fx.In

	Logger   log.Factory
	Channels map[string]*rabbitmq.Channel
}

func NewRabbitMQProducer(params RabbitMQProducerParams) *RabbitMQProducer {
	p := &RabbitMQProducer{
		Logger:   params.Logger,
		Channels: params.Channels,
	}
	// set the default channel and queue if exists in Channels and Queues
	if len(p.Channels) > 0 {
		for name, channel := range p.Channels {
			if name == "default" {
				p.Channel = channel
				break
			}
		}
	}

	return p
}

func (p *RabbitMQProducer) Publish(ctx context.Context, messageName string, message any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.Channel.Channel.PublishWithContext(
		ctx,
		info.ServiceName, // exchange
		messageName,      // routing key
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *RabbitMQProducer) SetChannel(channelName string) {
	channel, ok := p.Channels[channelName]
	if !ok {
		p.Logger.Bg().Error(fmt.Sprintf("channel %s not found", channelName))
	}
	p.Channel = channel
}
