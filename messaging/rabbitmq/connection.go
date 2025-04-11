package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/enesanbar/go-service/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Connection struct {
	logger        log.Factory
	Conn          *amqp.Connection
	ConnCloseChan chan *amqp.Error
	Config        *Config
}

type ConnectionParams struct {
	fx.In

	Logger log.Factory
	Config *Config
}

func NewConnector(p ConnectionParams) (*Connection, error) {
	return &Connection{
		logger: p.Logger,
		Config: p.Config,
	}, nil
}

func (c *Connection) connect() {
	c.logger.Bg().
		With(zap.String("db", c.Config.Host)).
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

		time.Sleep(5 * time.Second)
	}

	connChan := make(chan *amqp.Error)
	conn.NotifyClose(connChan)

	c.logger.Bg().
		With(zap.String("db", c.Config.Host)).
		With(zap.String("name", c.Config.Name)).
		Info("connected to rabbitmq")

	c.Conn = conn
	c.ConnCloseChan = connChan
}

func (c *Connection) Start(ctx context.Context) error {
	for {
		select {
		case err := <-c.ConnCloseChan:
			c.logger.Bg().
				With(zap.String("db", c.Config.Host)).
				With(zap.String("name", c.Config.Name)).
				With(zap.Error(err)).
				Error("connection closed")
			c.Close(ctx)
			c.connect()
		}

	}
	return nil
}

// Close closes the connection to RabbitMQ
func (c *Connection) Close(ctx context.Context) error {
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
