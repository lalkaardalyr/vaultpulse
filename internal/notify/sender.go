package notify

// Sender is the interface implemented by notification clients.
// Any type that can deliver a plain-text message satisfies this interface,
// making it straightforward to swap or mock notification backends in tests.
type Sender interface {
	// Send delivers the given message to the notification service.
	// It returns a non-nil error if delivery fails for any reason.
	Send(message string) error
}

// NoopSender is a Sender that silently discards every message.
// It is useful as a default when no notification backend is configured.
type NoopSender struct{}

// Send implements Sender. It always returns nil without doing anything.
func (n *NoopSender) Send(_ string) error { return nil }
