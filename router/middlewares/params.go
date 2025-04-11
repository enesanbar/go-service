package middlewares

import (
	"github.com/enesanbar/go-service/config"
	"github.com/enesanbar/go-service/log"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Env        string `name:"environment"`
	BaseConfig *config.Base
	Logger     log.Factory
}
