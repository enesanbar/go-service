package rabbitmq

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"github.com/enesanbar/go-service/core/messaging/consumer"
	"github.com/enesanbar/go-service/core/wiring"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ConsumersParams struct {
	fx.In

	Conf            config.Config
	Logger          log.Factory
	Queues          map[string]*Queue
	Channels        map[string]*Channel
	MessageHandlers map[string]consumer.MessageHandler
	Propagator      propagation.TextMapPropagator
	TracerProvider  *tracesdk.TracerProvider
}

func Consumers(p ConsumersParams) ([]wiring.Runnable, error) {
	runnables := make([]wiring.Runnable, 0)

	cfg := p.Conf.GetSliceOfObjects("rabbitmq.consumers")

	for _, v := range cfg {
		cfg, err := NewConsumerConfig(v)
		if err != nil {
			p.Logger.Bg().Error("failed to create consumer config")
			continue
		}

		channel, ok := p.Channels[cfg.Channel]
		if !ok {
			p.Logger.Bg().
				With(zap.String("channel", cfg.Channel)).
				Error("channel not found for exchange. please check the channel configuration in your configuration")
			continue
		}

		queue, ok := p.Queues[cfg.Queue]
		if !ok {
			p.Logger.Bg().
				With(zap.String("channel", cfg.Channel)).
				With(zap.String("queue", cfg.Queue)).
				Error("queue not found for exchange. please check the queue configuration in your configuration")
			continue
		}

		o := NewRabbitMQConsumer(ConsumerParams{
			Logger:          p.Logger,
			Config:          cfg,
			Channel:         channel,
			Queue:           queue,
			MessageHandlers: p.MessageHandlers,
			Propagator:      p.Propagator,
			TracerProvider:  p.TracerProvider,
		})
		runnables = append(runnables, o)
	}

	return runnables, nil
}
