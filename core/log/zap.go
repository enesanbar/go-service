package log

import (
	"github.com/enesanbar/go-service/core/info"
	"github.com/enesanbar/go-service/core/osutil"
	"go.uber.org/fx"

	"go.uber.org/zap"
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
		logger, err = zap.NewDevelopment()
		logger, err = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       true,
			Encoding:          "console",
			DisableStacktrace: true,
			EncoderConfig:     zap.NewProductionEncoderConfig(),
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},
		}.Build(
			zap.AddCallerSkip(1),
			zap.AddCaller(),
			zap.Fields(zap.String("version", info.Version)),
		)
	} else {
		//logger, err = zap.NewProduction(
		//	zap.AddStacktrace(zapcore.ErrorLevel),
		//	zap.AddCallerSkip(1),
		//	zap.AddCaller(),
		//	zap.Fields(zap.String("version", info.Version)),
		//)
		logger, err = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:          "json",
			DisableStacktrace: true,
			EncoderConfig:     zap.NewProductionEncoderConfig(),
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},
		}.Build(
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
