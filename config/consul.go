package config

import (
	"fmt"

	"github.com/enesanbar/go-service/info"
	"github.com/enesanbar/go-service/log"
	"github.com/enesanbar/go-service/osutil"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	EnvDev     = "dev"
	EnvTest    = "test"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

var consulHost string = osutil.GetEnv("CONSUL_HOST", "localhost")

type ConsulConfigProvider struct {
	viper  *viper.Viper
	logger log.Factory
	env    string
}

type Params struct {
	fx.In

	Logger log.Factory
	Viper  *viper.Viper
	Env    string `name:"environment"`
}

func NewConsulProvider(p Params) (*ConsulConfigProvider, error) {
	c := &ConsulConfigProvider{
		logger: p.Logger,
		env:    p.Env,
	}

	base, err := c.ReadRemoteProperties()
	if err != nil {
		return nil, err
	}

	err = p.Viper.MergeConfigMap(base)
	if err != nil {
		return nil, err
	}

	p.Viper.AutomaticEnv()
	c.viper = p.Viper

	return c, nil
}

func (c *ConsulConfigProvider) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *ConsulConfigProvider) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

func (c *ConsulConfigProvider) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *ConsulConfigProvider) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *ConsulConfigProvider) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *ConsulConfigProvider) ReadRemoteProperties() (map[string]interface{}, error) {
	c.logger.Bg().Info(
		"loading service configuration",
		zap.String("config_source", SourceConsul),
		zap.String("consul_address", consulHost),
	)

	extension := "yaml"
	paths := [...]string{
		"go-service",
		fmt.Sprintf("go-service_%s", c.env),
		info.ServiceName,
		fmt.Sprintf("%s_%s", info.ServiceName, c.env),
	}

	base := viper.New()
	base.SetConfigType(extension)

	for _, path := range paths {
		conf := viper.New()
		conf.SetConfigType("yaml")
		err := conf.AddRemoteProvider(
			"consul",
			fmt.Sprintf("http://%s:8500", consulHost),
			fmt.Sprintf("/go-config/%s.%s", path, extension))
		if err != nil {
			c.logger.Bg().Info("Unable to add consul provider.", zap.Error(err))
			return nil, err
		}

		c.logger.Bg().Info("reading config file", zap.String("filename", path))
		err = conf.ReadRemoteConfig()
		if err != nil {
			return nil, err
		}

		_ = base.MergeConfigMap(conf.AllSettings())
	}

	return base.AllSettings(), nil
}
