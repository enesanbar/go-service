package rabbitmq

import (
	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
	"go.uber.org/zap"
)

func RabbitMQQueues(conf config.Config, logger log.Factory, channels map[string]*Channel) (map[string]*Queue, error) {
	if len(channels) == 0 {
		return nil, nil
	}

	cfg := conf.GetStringMap("datasources.rabbitmq.queues")

	queues := make(map[string]*Queue)

	for queueName, v := range cfg {
		channelName := v.(map[string]interface{})[PropertyChannel].(string)
		channel, ok := channels[channelName]
		if !ok {
			logger.Bg().
				With(zap.String("queue", queueName)).
				With(zap.String("channel", channelName)).
				Error("channel not found for queue. please check the channel configuration in your configuration")
			continue
		}

		queue, err := NewQueue(QueueParams{
			Channel: channel,
			Logger:  logger,
			Config: &QueueConfig{
				Name:       queueName,
				Channel:    channel,
				Durable:    v.(map[string]interface{})[PropertyDurable].(bool),
				AutoDelete: v.(map[string]interface{})[PropertyAutoDelete].(bool),
				Exclusive:  v.(map[string]interface{})[PropertyExclusive].(bool),
				NoWait:     v.(map[string]interface{})[PropertyNoWait].(bool),
			},
		})

		if err != nil {
			logger.Bg().
				With(zap.String("queue", queueName)).
				With(zap.String("channel", channelName)).
				With(zap.Error(err)).
				Error("failed to create queue")
			panic(err)
		}

		queues[queueName] = queue
		logger.Bg().
			With(zap.String("queue", queueName)).
			With(zap.String("channel", channelName)).
			Info("queue created")
	}

	return queues, nil
}
