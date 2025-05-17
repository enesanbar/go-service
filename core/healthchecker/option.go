package healthchecker

type options struct {
	Probes []Probe
}

var defaultHealthCheckerOptions = options{
	Probes: []Probe{},
}

type Option func(o *options)

func WithProbes(p ...Probe) Option {
	return func(o *options) {
		o.Probes = p
	}
}
