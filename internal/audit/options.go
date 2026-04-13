package audit

import "io"

// Option configures a Logger.
type Option func(*Logger)

// WithWriter overrides the output writer used by the Logger.
func WithWriter(w io.Writer) Option {
	return func(l *Logger) {
		if w != nil {
			l.writer = w
		}
	}
}

// NewWithOptions constructs a Logger applying the provided options.
func NewWithOptions(opts ...Option) *Logger {
	l := New(nil)
	for _, o := range opts {
		o(l)
	}
	return l
}
