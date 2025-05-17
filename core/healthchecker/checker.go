package healthchecker

import "context"

type Result struct {
	Success       bool                    `json:"success"`
	ProbesResults map[string]*ProbeResult `json:"probes"`
}

// HealthChecker runs a set of health checks provided by the application developer and returns the results.
// If any of the checks fail, the overall result is considered a failure.
type HealthChecker struct {
	probes []Probe
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		probes: []Probe{},
	}
}

func (c *HealthChecker) AddProbe(p Probe) *HealthChecker {
	c.probes = append(c.probes, p)

	return c
}

func (c *HealthChecker) Run(ctx context.Context) *Result {

	success := true
	probeResults := map[string]*ProbeResult{}

	for _, p := range c.probes {

		pr := p.Check(ctx)

		success = success && pr.Success
		probeResults[p.Name()] = pr
	}

	return &Result{
		Success:       success,
		ProbesResults: probeResults,
	}
}
