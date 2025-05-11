package log

import (
	"github.com/enesanbar/go-service/core/osutil"

	"go.uber.org/fx"

	"github.com/enesanbar/go-service/core/info"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Params struct {
	fx.In

	Env string `name:"environment"`
}

// NewZapLogger constructs a new logger.
func NewZapLogger() (*zap.Logger, error) {
	env := osutil.GetEnv("DEPLOY_TYPE", "dev")

	var logger *zap.Logger
	var err error
	if env != "prod" {
		logger, err = zap.NewDevelopment(
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.AddCallerSkip(1),
		)
	} else {
		logger, err = zap.NewProduction(
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.AddCallerSkip(1),
			zap.AddCaller(),
			zap.Fields(zap.String("version", info.Version)),
		)
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}
