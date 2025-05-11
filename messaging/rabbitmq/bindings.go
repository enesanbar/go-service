package rabbitmq

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RabbitMQBindingsParams struct {
	fx.In

	Conf      config.Config
	Logger    log.Factory
	Queues    map[string]*Queue    `optinal:"true"`
	Exchanges map[string]*Exchange `optinal:"true"`
}

func RabbitMQBindings(p RabbitMQBindingsParams) error {
	if len(p.Queues) == 0 || len(p.Exchanges) == 0 {
		return nil
	}

	cfg := p.Conf.GetSliceOfObjects("datasources.rabbitmq.bindings")

	for _, v := range cfg {
		exchangeName := v.(map[string]interface{})[PropertyExchange].(string)
		queueName := v.(map[string]interface{})[PropertyQueue].(string)

		exchange, ok := p.Exchanges[exchangeName]
		if !ok {
			p.Logger.Bg().
				With(zap.String("exchange", exchangeName)).
				Error("exchange not found for binding. please check the exchange configuration in your configuration")
			continue
		}
		queue, ok := p.Queues[queueName]
		if !ok {
			p.Logger.Bg().
				With(zap.String("queue", queueName)).
				Error("queue not found for binding. please check the queue configuration in your configuration")
			continue
		}

		routingKeys := v.(map[string]interface{})[PropertyRoutingKeys].([]interface{})

		for _, routingKey := range routingKeys {
			err := queue.Channel.Channel.QueueBind(
				queue.Config.Name,
				routingKey.(string),
				exchange.Config.Name,
				v.(map[string]interface{})[PropertyNoWait].(bool),
				nil,
			)

			if err != nil {
				p.Logger.Bg().
					With(zap.String("exchange", exchangeName)).
					With(zap.String("queue", queueName)).
					With(zap.String("routingKey", routingKey.(string))).
					With(zap.Error(err)).
					Error("failed to create binding")
				return err
			}

			p.Logger.Bg().
				With(zap.String("binding", routingKey.(string))).
				With(zap.String("exchange", exchangeName)).
				With(zap.String("queue", queueName)).
				Info("binding created")
		}

	}

	return nil
}
