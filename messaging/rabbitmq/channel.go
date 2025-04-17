package rabbitmq

import (
	"context"
	"time"

	"github.com/enesanbar/go-service/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Channel struct {
	logger        log.Factory
	Channel       *amqp.Channel
	ConnCloseChan chan *amqp.Error
	Config        *ChannelConfig
}

type ChannelParams struct {
	fx.In

	Logger log.Factory
	Config *ChannelConfig
}

func NewChannelConnector(p ChannelParams) (*Channel, error) {
	return &Channel{
		logger: p.Logger,
		Config: p.Config,
	}, nil
}

func (c *Channel) connect() {
	// TODO: Fix the data race (access and write to Conn) between the channel and the connection on restart
	channel, err := c.Config.Connection.Conn.Channel()
	if err != nil {
		c.logger.Bg().
			With(zap.String("db", c.Config.Connection.Config.Host)).
			With(zap.String("name", c.Config.Name)).
			With(zap.Error(err)).
			Error("failed to create channel, retrying in 5 seconds")
		return
	}

	c.logger.Bg().
		With(zap.String("name", c.Config.Name)).
		With(zap.String("connection", c.Config.Connection.Name())).
		Info("created channel to rabbitmq")

	c.Channel = channel
	c.ConnCloseChan = make(chan *amqp.Error)
	c.Channel.NotifyClose(c.ConnCloseChan)
}

func (c *Channel) Start(ctx context.Context) error {
	c.logger.Bg().
		With(zap.String("name", c.Config.Name)).
		Info("starting channel watcher")
	for {
		if c.ConnCloseChan == nil {
			time.Sleep(5 * time.Second)
		}

		// if c.Config.Connection.Conn.IsClosed() {
		// 	c.logger.Bg().
		// 		With(zap.String("name", c.Config.Name)).
		// 		With(zap.String("connection", c.Config.Connection.Config.Host)).
		// 		Error("connection to RabbitMQ is not open yet, waiting for it to open")
		// 	time.Sleep(5 * time.Second)
		// 	continue
		// }

		select {
		case err := <-c.ConnCloseChan:
			c.logger.Bg().
				With(zap.String("name", c.Config.Name)).
				With(zap.String("db", c.Config.Connection.Name())).
				With(zap.Error(err)).
				Error("Channel closed, reconnecting")
			if !c.Channel.IsClosed() {
				c.Close(ctx)
			}
			c.connect()
			time.Sleep(5 * time.Second)
		}

	}
	return nil
}

// Close closes the Channel to RabbitMQ
func (c *Channel) Close(ctx context.Context) error {
	err := c.Channel.Close()
	if err != nil {
		return err
	}

	c.logger.Bg().
		With(zap.String("channel", c.Config.Name)).
		With(zap.String("connection", c.Config.Connection.Name())).
		Info("closed Channel to rabbitmq")

	return nil
}

func (c *Channel) Name() string {
	return c.Config.Name
}
