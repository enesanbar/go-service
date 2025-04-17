package rabbitmq

type ExchangeConfig struct {
	Name       string
	Type       string
	Channel    *Channel
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}
