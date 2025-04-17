package rabbitmq

type QueueConfig struct {
	Name       string
	Channel    *Channel
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}
