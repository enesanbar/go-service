package rabbitmq

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/zap"
)

// Connections return a map of connections configured in the configuration file.
func Connections(conf config.Config, logger log.Factory) (map[string]*Connection, error) {
	prefix := "rabbitmq.connections"
	cfg := conf.GetStringMap(prefix)
	connections := make(map[string]*Connection)
	for k := range cfg {
		config, err := NewConnectionConfig(conf, k)
		if err != nil {
			logger.Bg().
				With(zap.String("connection", k)).
				With(zap.Error(err)).
				Error("failed to create connection config")
			return nil, err
		}
		conn := &Connection{
			logger:        logger,
			Config:        config,
			AppStopSignal: make(chan struct{}),
		}
		err = conn.connect()
		if err != nil {
			return nil, err
		}
		connections[k] = conn
	}

	return connections, nil
}
