package config

import (
	"github.com/enesanbar/go-service/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	BaseConfig       = "base"
	DefaultExtension = "yaml"
)

type FileConfigProviderParams struct {
	fx.In

	Logger log.Factory
	Viper  *viper.Viper
	Env    string `name:"environment"`
}

type FileConfigProvider struct {
	logger log.Factory
	viper  *viper.Viper
	env    string
}

func NewFileConfigProvider(p FileConfigProviderParams) (*FileConfigProvider, error) {
	f := &FileConfigProvider{logger: p.Logger, viper: p.Viper, env: p.Env}

	properties, err := f.ReadLocalProperties()
	if err != nil {
		return nil, err
	}

	if err := f.viper.MergeConfigMap(properties); err != nil {
		return nil, err
	}

	return f, nil
}

func (c *FileConfigProvider) ReadLocalProperties() (map[string]interface{}, error) {
	c.logger.Bg().Info(
		"loading service configuration",
		zap.String("config_source", SourceFile),
	)

	configDir := "./config"

	paths := [...]string{BaseConfig, c.env}

	base := viper.New()
	base.SetConfigType(DefaultExtension)

	for _, path := range paths {
		conf := viper.New()
		conf.SetConfigType(DefaultExtension)
		conf.SetConfigName(path)
		conf.AddConfigPath(configDir)
		err := conf.ReadInConfig()
		if err != nil {
			// c.logger.Bg().With(
			// 	zap.String("config_path", configDir),
			// 	zap.String("config_name", path),
			// 	zap.String("config_type", DefaultExtension),
			// ).Error("failed to read config file", zap.Error(err))
			continue
		}
		_ = base.MergeConfigMap(conf.AllSettings())
	}

	return base.AllSettings(), nil
}

func (c *FileConfigProvider) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *FileConfigProvider) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

func (c *FileConfigProvider) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *FileConfigProvider) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *FileConfigProvider) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}
