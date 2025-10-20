# Registration with the application

1. Implement Probe interface to register with the application.
2. Provide the constructor
3. Define more health checkers in the same way to check different services.

Implementation of the Probe interface:
```go
type Checker struct {
}

type Params struct {
    fx.In
    Repository SomeRepository
}

func NewChecker(p Params) *Checker {
    return &Checker{}
}

// Name returns the name of the health checker.
func (c *Checker) Name() string {
    return "my-health-checker"
}

// Check performs the health check and returns the result.
func (c *Checker) Check(ctx context.Context) *healthchecker.ProbeResult {
    // Perform your health check logic here
    // For example, check if a database connection is healthy
    err := c.Repository.Ping()    
    if err != nil {
        return healthchecker.NewProbeResult(false, "database ping failed: "+err.Error())
    }
    return healthchecker.NewProbeResult(true, "database ping success")
}
```

Registration of the health checker:
```go
package health

import (
	"github.com/enesanbar/go-service/core/healthchecker"
	"go.uber.org/fx"
)

var Module = fx.Options(
	factories,
)

var factories = fx.Provide(
	healthchecker.AsHealthCheckerProbe(NewChecker),
)

```