package instrumentation

import (
	"github.com/enesanbar/go-service/info"
	"github.com/enesanbar/go-service/osutil"
	instana "github.com/instana/go-sensor"
)

func NewInstanaSensor() *instana.Sensor {
	if osutil.GetEnv("INSTANA_ENABLED", "false") == "false" {
		return nil
	}
	return instana.NewSensor(info.ServiceName)
}
