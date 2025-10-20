package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/messaging/messages"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/fx"
)

type Producer struct {
	Logger     log.Factory
	Connection *Connection
	Propagator propagation.TextMapPropagator
}

type ProducerParams struct {
	fx.In

	Logger      log.Factory
	Connections map[string]*Connection
	Propagator  propagation.TextMapPropagator
}

// NewRabbitMQProducer creates a pointer to the new instance of the Producer
func NewRabbitMQProducer(params ProducerParams) (*Producer, error) {
	p := &Producer{
		Logger:     params.Logger,
		Propagator: params.Propagator,
	}

	if len(params.Connections) == 0 {
		return nil, fmt.Errorf("no connections found. please check the connection configuration in your configuration")
	}

	// Use the first connection for the producer
	for _, conn := range params.Connections {
		if conn == nil {
			continue
		}
		p.Connection = conn
		p.Logger.Bg().Info("using connection for producer", zap.String("connection", conn.Name()))
		break
	}

	if p.Connection == nil {
		return nil, fmt.Errorf("no valid connection found for producer")
	}

	return p, nil
}

func (p *Producer) Publish(ctx context.Context, messageName string, payload any) error {
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

	channel, err := p.Connection.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer func() {
		if err := channel.Close(); err != nil {
			p.Logger.Bg().Error("failed to close channel", zap.Error(err))
		}
	}()

	// TODO: Optionally log the message
	return channel.PublishWithContext(
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

func (p *Producer) enrichMessageWithTrace(ctx context.Context, message *messages.Message[any]) {
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
