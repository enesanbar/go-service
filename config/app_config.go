package config

type Base struct {
	Environment string `json:"environment" yaml:"environment"`
	Debug       bool   `json:"debug" yaml:"debug"`
}

func NewBaseConfig(config Config) *Base {
	return &Base{
		Environment: config.GetString("env"),
		Debug:       config.GetBool("debug"),
	}
}

func (b *Base) IsVerbose() bool {
	return b.Environment == EnvDev || b.Debug
}
