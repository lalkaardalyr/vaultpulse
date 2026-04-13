package runner

// Option is a functional option for configuring a Runner.
type Option func(*Runner)

// WithNotifierThreshold is a placeholder demonstrating how additional
// runner-level options can be layered without changing the New signature.
//
// Example:
//
//	r, err := runner.New(cfg, w, runner.WithNotifierThreshold(alerts.Warning))
func WithNotifierThreshold(_ interface{}) Option {
	return func(_ *Runner) {
		// Reserved for future use: pass threshold down to the notifier.
	}
}

// applyOptions applies all provided Option functions to the Runner.
func applyOptions(r *Runner, opts []Option) {
	for _, o := range opts {
		o(r)
	}
}
