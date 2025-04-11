package validation

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Options(
	factories,
)

var factories = fx.Provide(
	validator.New,
	fx.Annotated{
		Name:   "go_playground",
		Target: NewGoPlayground,
	},
)
