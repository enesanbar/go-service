package mysql

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/zap"
)

func MySQLConnections(conf config.Config, logger log.Factory) (map[string]*Connection, error) {
	prefix := "datasources.mysql"
	cfg := conf.GetStringMap(prefix)
	connections := make(map[string]*Connection)
	for k := range cfg {
		config, err := NewConfig(conf, k)
		if err != nil {
			logger.Bg().
				With(zap.String("connection", k)).
				With(zap.Error(err)).
				Error("failed to create connection config")
			return nil, err
		}
		conn := NewConnection(ConnectionParams{
			Config: config,
			Logger: logger,
		})
		// TODO: this does nothing. Implement lazy start later
		err = conn.Start()
		if err != nil {
			return nil, err
		}
		connections[k] = conn
	}

	return connections, nil
}
