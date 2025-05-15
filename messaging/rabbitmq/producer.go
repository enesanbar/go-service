package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/messaging/messages"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/fx"
)

type RabbitMQProducer struct {
	Logger     log.Factory
	Channel    *Channel
	Channels   map[string]*Channel
	Propagator propagation.TextMapPropagator
}

type RabbitMQProducerParams struct {
	fx.In

	Logger     log.Factory
	Channels   map[string]*Channel
	Propagator propagation.TextMapPropagator
}

func NewRabbitMQProducer(params RabbitMQProducerParams) *RabbitMQProducer {
	p := &RabbitMQProducer{
		Logger:     params.Logger,
		Channels:   params.Channels,
		Propagator: params.Propagator,
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

func (p *RabbitMQProducer) Publish(ctx context.Context, messageName string, payload any) error {
	message := messages.Message[any]{
		Metadata: messages.Metadata{
			PublisherName: info.ServiceName,
			PublishDate:   time.Now().UTC(),
			MessageName:   messageName,
		},
		Payload: payload,
	}

	// Enrich the message with trace information
	p.enrichMessageWithTrace(ctx, &message)

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

func (p *RabbitMQProducer) enrichMessageWithTrace(ctx context.Context, message *messages.Message[any]) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	// Properly inject trace context using OpenTelemetry's propagator
	carrier := propagation.MapCarrier{}
	p.Propagator.Inject(ctx, carrier)

	message.Metadata.Traceparent = carrier["traceparent"]
	message.Metadata.Tracestate = carrier["tracestate"]

	// This is the CURRENT span ID (the one sending the message)
	message.Metadata.SpanID = span.SpanContext().SpanID().String()
}
