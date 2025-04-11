package mongodb

import (
	"context"
	"fmt"

	"github.com/instana/go-sensor/instrumentation/instamongo"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/enesanbar/go-service/log"
	instana "github.com/instana/go-sensor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Connector struct {
	logger log.Factory
	sensor *instana.Sensor
}

type ConnectionParams struct {
	fx.In

	Sensor *instana.Sensor `optional:"true"`
	Logger log.Factory
}

func NewConnector(p ConnectionParams) (*Connector, error) {
	return &Connector{
		logger: p.Logger,
		sensor: p.Sensor,
	}, nil
}

func (c *Connector) Connect(cfg *Config) (*mongo.Client, error) {
	c.logger.Bg().
		With(zap.String("db", cfg.Host)).
		Info("connecting to mongo database")

	//DSN := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	DSN := fmt.Sprintf("mongodb://%s", cfg.Host)
	clientOptions := &options.ClientOptions{
		ServerSelectionTimeout: &cfg.Timeout,
		SocketTimeout:          &cfg.Timeout,
		ConnectTimeout:         &cfg.Timeout,
		MaxPoolSize:            &cfg.MaxPoolSize,
		MinPoolSize:            &cfg.MinPoolSize,
		MaxConnIdleTime:        &cfg.MaxConnectionIdletime,
	}

	if cfg.ReplicaSetName != "" {
		clientOptions.ReplicaSet = &cfg.ReplicaSetName
	}

	if cfg.User != "" && cfg.Pass != "" {
		clientOptions.Auth = &options.Credential{
			AuthSource: cfg.AuthDB,
			Username:   cfg.User,
			Password:   cfg.Pass,
		}
	}

	var client *mongo.Client
	var err error

	if c.sensor != nil {
		client, err = instamongo.NewClient(
			c.sensor,
			options.Client().ApplyURI(DSN),
			clientOptions,
		)
	} else {
		client, err = mongo.NewClient(options.Client().ApplyURI(DSN), clientOptions)
	}

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Connector) Start(ctx context.Context) error {
	return nil
}

func (c *Connector) Name(ctx context.Context) string {
	return "mongodb"
}
