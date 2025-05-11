package healthchecker

import (
	"context"
)

// HealthCheckerProbe is an interface that defines a health check probe.
// It should be implemented by applications that want to perform health checks and
// provided to fx as a dependency with AsHealthCheckerProbe(NewChecker),
type HealthCheckerProbe interface {
	Name() string
	Check(ctx context.Context) *HealthCheckerProbeResult
}

type HealthCheckerProbeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHealthCheckerProbeResult(success bool, message string) *HealthCheckerProbeResult {
	return &HealthCheckerProbeResult{
		Success: success,
		Message: message,
	}
}
