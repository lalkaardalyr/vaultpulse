package alerts

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Notifier sends alerts to one or more destinations.
type Notifier struct {
	writer            io.Writer
	warnThreshold     time.Duration
	criticalThreshold time.Duration
}

// NotifierOption configures a Notifier.
type NotifierOption func(*Notifier)

// WithWriter overrides the default stdout writer.
func WithWriter(w io.Writer) NotifierOption {
	return func(n *Notifier) {
		n.writer = w
	}
}

// NewNotifier creates a Notifier with the given thresholds.
func NewNotifier(warnThreshold, criticalThreshold time.Duration, opts ...NotifierOption) *Notifier {
	n := &Notifier{
		writer:            os.Stdout,
		warnThreshold:     warnThreshold,
		criticalThreshold: criticalThreshold,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// Notify evaluates the expiry time for a secret and writes an alert if warranted.
func (n *Notifier) Notify(secretPath string, expiresAt time.Time) error {
	alert := NewAlert(secretPath, expiresAt, n.warnThreshold, n.criticalThreshold)
	if alert.Level == LevelInfo {
		return nil
	}
	_, err := fmt.Fprintln(n.writer, alert.Message)
	return err
}

// NotifyAll sends alerts for a map of secret paths to expiry times.
func (n *Notifier) NotifyAll(secrets map[string]time.Time) error {
	for path, expiry := range secrets {
		if err := n.Notify(path, expiry); err != nil {
			return fmt.Errorf("notifier: failed to notify for %s: %w", path, err)
		}
	}
	return nil
}
