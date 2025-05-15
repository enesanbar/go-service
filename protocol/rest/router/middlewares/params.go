package middlewares

import (
	"github.com/enesanbar/go-service/core/config"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Env        string `name:"environment"`
	BaseConfig *config.Base
	Logger     log.Factory
}
