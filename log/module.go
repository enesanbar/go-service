package log

import (
	"go.uber.org/fx"
)

// Module is the loggerfx module that can be passed into an Fx app.
var Module = fx.Options(
	factories,
)

var factories = fx.Provide(
	NewZapLogger,
	NewFactory,
)
