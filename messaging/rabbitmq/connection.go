package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/enesanbar/go-service/core/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Connection struct {
	logger        log.Factory
	Conn          *amqp.Connection
	ConnCloseChan chan *amqp.Error
	Config        *ConnectionConfig
	AppStopSignal chan struct{}
}

type ConnectionParams struct {
	fx.In

	Logger log.Factory
	Config *ConnectionConfig
}

func NewConnector(p ConnectionParams) (*Connection, error) {
	return &Connection{
		logger:        p.Logger,
		Config:        p.Config,
		AppStopSignal: make(chan struct{}),
	}, nil
}

func (c *Connection) connect() error {
	c.logger.Bg().
		With(zap.String("host", c.Config.Host)).
		With(zap.String("name", c.Config.Name)).
		Info("connecting to rabbitmq")

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.Config.User,
		c.Config.Pass,
		c.Config.Host,
		c.Config.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		c.logger.Bg().
			With(zap.String("db", c.Config.Host)).
			With(zap.String("name", c.Config.Name)).
			With(zap.Error(err)).
			Error("failed to connect to rabbitmq, retrying in 5 seconds")
		return err
	}

	connChan := make(chan *amqp.Error)
	conn.NotifyClose(connChan)

	c.logger.Bg().
		With(zap.String("host", c.Config.Host)).
		With(zap.String("name", c.Config.Name)).
		Info("connected to rabbitmq")

	// TODO: Fix the data race (access and write to Conn) between the channel and the connection on restart
	c.Conn = conn
	c.ConnCloseChan = connChan

	return nil
}

func (c *Connection) Start(ctx context.Context) error {
	c.logger.Bg().
		With(zap.String("name", c.Config.Name)).
		Info("starting connection watcher for rabbitmq connection")

	for {
		if c.ConnCloseChan == nil {
			time.Sleep(5 * time.Second)
		}
		select {
		case err := <-c.ConnCloseChan:
			c.logger.Bg().
				With(zap.String("host", c.Config.Host)).
				With(zap.String("name", c.Config.Name)).
				With(zap.Error(err)).
				Error("connection closed")
			if !c.Conn.IsClosed() {
				return c.Close(ctx)
			}
			c.connect()
			time.Sleep(5 * time.Second)
		case <-c.AppStopSignal:
			c.logger.For(ctx).Info("context done, stopping the connection watcher")
			return nil
		}

	}
}

// Close closes the connection to RabbitMQ
func (c *Connection) Close(ctx context.Context) error {
	c.AppStopSignal <- struct{}{}

	if c.Conn.IsClosed() {
		c.logger.Bg().
			With(zap.String("host", c.Config.Host)).
			With(zap.String("name", c.Config.Name)).
			Info("connection already closed")
		return nil
	}
	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.logger.Bg().
		With(zap.String("db", c.Config.Host)).
		With(zap.String("name", c.Config.Name)).
		Info("closed connection to rabbitmq")

	return nil
}

func (c *Connection) Name() string {
	return c.Config.Name
}
