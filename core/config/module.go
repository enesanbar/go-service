package config

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	factories,
)

var ConsulModule = fx.Provide(
	NewConsulProvider,
	func(consul *ConsulConfigProvider) Config { return consul },
)

var FileModule = fx.Provide(
	NewFileConfigProvider,
	func(file *FileConfigProvider) Config { return file },
)

var factories = fx.Options(
	NewConfig(),
	fx.Provide(
		NewViper,
		NewBaseConfig,
		fx.Annotated{
			Name:   "environment",
			Target: DetermineEnvironment,
		},
	),
)
