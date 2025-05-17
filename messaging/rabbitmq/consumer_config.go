package rabbitmq

type ConsumerConfig struct {
	ConsumerTag   string
	Channel       string
	Queue         string
	AutoAck       bool
	Exclusive     bool
	NoLocal       bool
	NoWait        bool
	Requeue       bool
	PrefetchCount int
}

func NewConsumerConfig(cfg interface{}) (*ConsumerConfig, error) {
	config := cfg.(map[string]interface{})

	consumerTag, ok := config[PropertyConsumerTag]
	if !ok {
		consumerTag = "" // library generates a random tag
	}

	queueName, ok := config[PropertyQueue]
	if !ok {
		panic("queue name is required")
	}

	channelName, ok := config[PropertyChannel]
	if !ok {
		panic("channel name is required")
	}

	autoAck, ok := config[PropertyAutoAck]
	if !ok {
		autoAck = true
	}

	exclusive, ok := config[PropertyExclusive]
	if !ok {
		exclusive = false
	}

	noLocal, ok := config[PropertyNoLocal]
	if !ok {
		noLocal = false
	}

	noWait, ok := config[PropertyNoWait]
	if !ok {
		noWait = false
	}

	return &ConsumerConfig{
		ConsumerTag: consumerTag.(string),
		Channel:     channelName.(string),
		Queue:       queueName.(string),
		AutoAck:     autoAck.(bool),
		Exclusive:   exclusive.(bool),
		NoLocal:     noLocal.(bool),
		NoWait:      noWait.(bool),
	}, nil
}
