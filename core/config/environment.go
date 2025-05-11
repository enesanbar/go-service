package config

import (
	"fmt"
	"os"

	"github.com/enesanbar/go-service/core/osutil"
	"github.com/enesanbar/go-service/core/log"
	"go.uber.org/zap"
)

var envList = [...]string{EnvDev, EnvTest, EnvStaging, EnvProd}

func DetermineEnvironment(log log.Factory) string {
	env := osutil.GetEnv("DEPLOY_TYPE", EnvDev)

	err := validateEnvironment(env)
	if err != nil {
		log.Bg().Error("Invalid DEPLOY_TYPE", zap.Error(err))
		os.Exit(1)
	}

	log.Bg().Info("environment is activated", zap.String("env", env))
	return env
}

func validateEnvironment(env string) error {
	for _, e := range envList {
		if env == e {
			return nil
		}
	}
	return fmt.Errorf("environment is not supported. supported environments: %s", envList)
}
