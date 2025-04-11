package config

import (
	"github.com/enesanbar/go-service/osutil"
	"go.uber.org/fx"
)

const (
	SourceKey    = "CONFIG_SOURCE"
	SourceFile   = "file"
	SourceConsul = "consul"
)

type Config interface {
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetInt(key string) int
	GetBool(key string) bool
	GetStringSlice(key string) []string
}

type CloudProvider interface {
	Config
	ReadRemoteProperties() (map[string]interface{}, error)
}

type FileProvider interface {
	Config
	ReadLocalProperties() (map[string]interface{}, error)
}

func NewConfig() fx.Option {
	configSource := osutil.GetEnv(SourceKey, SourceFile)

	if configSource == SourceConsul {
		return ConsulModule
	}

	return FileModule
}
