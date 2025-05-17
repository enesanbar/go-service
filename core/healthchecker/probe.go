package healthchecker

import (
	"context"
)

// Probe is an interface that defines a health check probe.
// It should be implemented by applications that want to perform health checks and
// provided to fx as a dependency with AsHealthCheckerProbe(NewChecker),
type Probe interface {
	Name() string
	Check(ctx context.Context) *ProbeResult
}

type ProbeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewProbeResult(success bool, message string) *ProbeResult {
	return &ProbeResult{
		Success: success,
		Message: message,
	}
}
