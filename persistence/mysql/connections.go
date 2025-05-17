package mysql

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/zap"
)

// Connections return a map of connections configured in the configuration file.
func Connections(conf config.Config, logger log.Factory) (map[string]*Connection, error) {
	prefix := "mysql"
	cfg := conf.GetStringMap(prefix)
	connections := make(map[string]*Connection)
	for k := range cfg {
		currentConfig, err := NewConfig(conf, k)
		if err != nil {
			logger.Bg().
				With(zap.String("connection", k)).
				With(zap.Error(err)).
				Error("failed to create connection config")
			return nil, err
		}
		conn := NewConnection(ConnectionParams{
			Config: currentConfig,
			Logger: logger,
		})
		connections[k] = conn
	}

	return connections, nil
}
