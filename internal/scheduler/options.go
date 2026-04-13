package scheduler

import "time"

// Option is a functional option for configuring a Scheduler.
type Option func(*Scheduler)

// WithInterval overrides the polling interval of a Scheduler.
func WithInterval(d time.Duration) Option {
	return func(s *Scheduler) {
		if d > 0 {
			s.interval = d
		}
	}
}

// Apply applies a slice of options to the Scheduler.
func (s *Scheduler) Apply(opts ...Option) {
	for _, o := range opts {
		o(s)
	}
}
