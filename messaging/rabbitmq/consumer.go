package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"strings"

	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/messaging/consumer"
	"github.com/enesanbar/go-service/core/messaging/messages"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueueConsumer struct {
	logger          log.Factory
	Config          *ConsumerConfig
	Channel         *Channel
	Queue           *Queue
	MessageHandlers map[string]consumer.MessageHandler
	Propagator      propagation.TextMapPropagator
	Tracer          trace.Tracer
}

type ConsumerParams struct {
	Logger          log.Factory
	Config          *ConsumerConfig
	Channel         *Channel
	Queue           *Queue
	MessageHandlers map[string]consumer.MessageHandler
	Propagator      propagation.TextMapPropagator
	TracerProvider  *tracesdk.TracerProvider
}

// NewRabbitMQConsumer creates a pointer to the new instance of the QueueConsumer
// and a runnable group that can be used to start and stop the consumer
func NewRabbitMQConsumer(p ConsumerParams) *QueueConsumer {
	return &QueueConsumer{
		logger:          p.Logger,
		Config:          p.Config,
		Channel:         p.Channel,
		Queue:           p.Queue,
		MessageHandlers: p.MessageHandlers,
		Propagator:      p.Propagator,
		Tracer:          p.TracerProvider.Tracer(fmt.Sprintf("consumer-%s", p.Queue.Config.Name)),
	}
}

func (h *QueueConsumer) Start(ctx context.Context) error {
	if h.Channel == nil {
		return fmt.Errorf("channel is not set, check your configuration")
	}
	if h.Queue == nil {
		return fmt.Errorf("queue is not set, check your configuration")
	}

	// add recovery logic to the channel when channel/connection is closed
	msgs, err := h.Channel.Channel.Consume(
		h.Queue.Queue.Name,
		h.Config.ConsumerTag,
		h.Config.AutoAck,
		h.Config.Exclusive,
		h.Config.NoLocal,
		h.Config.NoWait,
		nil, // args
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
					h.logger.For(ctx).With(zap.Error(err)).Error("Failed to unmarshal message")
					return
				}

				key := fmt.Sprintf("%s-%s", h.Queue.Config.Name, message.Metadata.MessageName)
				handler, ok := h.MessageHandlers[key]
				if !ok {
					h.logger.For(ctx).With(zap.String("messageName", message.Metadata.MessageName)).Error("no handler found for message")
					return
				}

				// Unmarshal the payload to the correct type
				payload := handler.GetMessageType()
				message.UnmarshalPayload(payload)
				message.Payload = payload

				// Extract parent context from traceparent
				carrier := propagation.MapCarrier{
					"traceparent": message.Metadata.Traceparent,
					"tracestate":  message.Metadata.Tracestate,
				}
				ctx := h.Propagator.Extract(ctx, carrier)

				// Start a new span that:
				// 1. Continues the trace from traceparent
				// 2. Links to the span that sent the message (message.Metadata.SpanID)
				ctx, span := h.Tracer.Start(
					ctx,
					"processing: "+message.Metadata.MessageName,
					trace.WithSpanKind(trace.SpanKindConsumer),
					trace.WithLinks(trace.Link{
						SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
							TraceID:    parseTraceID(message.Metadata.Traceparent), // From traceparent
							SpanID:     parseSpanID(message.Metadata.SpanID),       // From metadata
							TraceFlags: trace.FlagsSampled,
							Remote:     true,
						}),
					}),
				)
				defer span.End()

				err = handler.Handle(ctx, message)
				if err != nil {
					h.logger.For(ctx).With(zap.Error(err)).Error("failed to handle message")
				}
			}(d)
		}
		h.logger.Bg().Info("RabbitMQ consumer stopped")
	}()
	h.logger.Bg().Info(fmt.Sprintf("RabbitMQ consumer started for queue %s", h.Queue.Config.Name))
	return nil
}

func (h *QueueConsumer) Stop(ctx context.Context) error {
	return nil
}

func parseTraceID(traceparent string) trace.TraceID {
	parts := strings.Split(traceparent, "-")
	if len(parts) >= 2 {
		tid, _ := trace.TraceIDFromHex(parts[1])
		return tid
	}
	return trace.TraceID{}
}

func parseSpanID(spanID string) trace.SpanID {
	sid, _ := trace.SpanIDFromHex(spanID)
	return sid
}
