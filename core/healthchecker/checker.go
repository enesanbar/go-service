package healthchecker

import "context"

type HealthCheckerResult struct {
	Success       bool                                 `json:"success"`
	ProbesResults map[string]*HealthCheckerProbeResult `json:"probes"`
}

// HealthChecker runs a set of health checks provided by the application developer and returns the results.
// If any of the checks fail, the overall result is considered a failure.
type HealthChecker struct {
	probes []HealthCheckerProbe
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		probes: []HealthCheckerProbe{},
	}
}

func (c *HealthChecker) AddProbe(p HealthCheckerProbe) *HealthChecker {
	c.probes = append(c.probes, p)

	return c
}

func (c *HealthChecker) Run(ctx context.Context) *HealthCheckerResult {

	success := true
	probeResults := map[string]*HealthCheckerProbeResult{}

	for _, p := range c.probes {

		pr := p.Check(ctx)

		success = success && pr.Success
		probeResults[p.Name()] = pr
	}

	return &HealthCheckerResult{
		Success:       success,
		ProbesResults: probeResults,
	}
}
