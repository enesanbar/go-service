package log

import (
	"os"

	"github.com/enesanbar/go-service/osutil"

	"go.uber.org/fx"

	"github.com/enesanbar/go-service/info"

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
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	logger = logger.With(
		zap.String("service", info.ServiceName),
		zap.String("hostname", hostname))

	return logger, nil
}
