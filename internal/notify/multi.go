package notify

import (
	"errors"
	"fmt"
)

// MultiSender fans out a single Send call to multiple Sender implementations.
// All senders are attempted; accumulated errors are returned as a combined error.
type MultiSender struct {
	senders []Sender
}

// NewMultiSender creates a MultiSender from the provided senders.
// Returns an error if no senders are supplied.
func NewMultiSender(senders ...Sender) (*MultiSender, error) {
	if len(senders) == 0 {
		return nil, fmt.Errorf("multi: at least one sender is required")
	}
	return &MultiSender{senders: senders}, nil
}

// Send delivers message to every registered sender.
// Errors from individual senders are collected and joined.
func (m *MultiSender) Send(message string) error {
	var errs []error
	for _, s := range m.senders {
		if err := s.Send(message); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
