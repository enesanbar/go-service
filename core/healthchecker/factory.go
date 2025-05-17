package healthchecker

type Factory interface {
	Create(options ...Option) (*HealthChecker, error)
}

type DefaultFactory struct{}

func NewDefaultFactory() Factory {
	return &DefaultFactory{}
}

func (f *DefaultFactory) Create(options ...Option) (*HealthChecker, error) {

	appliedOpts := defaultHealthCheckerOptions
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	checker := NewHealthChecker()

	for _, probe := range appliedOpts.Probes {
		checker.AddProbe(probe)
	}

	return checker, nil
}
