package rabbitmq

import (
	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
)

func RabbitMQConnections(conf config.Config, logger log.Factory) (map[string]*Connection, error) {
	cfg := conf.GetStringMap("datasources.rabbitmq")

	connections := make(map[string]*Connection)
	for k, v := range cfg {
		conn := &Connection{
			logger: logger,
			Config: &Config{
				Name: k,
				Host: v.(map[string]interface{})["host"].(string),
				Port: v.(map[string]interface{})["port"].(string),
				User: v.(map[string]interface{})["username"].(string),
				Pass: v.(map[string]interface{})["password"].(string),
			},
		}
		conn.connect()
		connections[k] = conn
	}

	return connections, nil
}
