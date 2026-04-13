package notify

// Message represents a notification payload sent by any Sender implementation.
type Message struct {
	// Body is the human-readable alert text.
	Body string

	// Severity indicates the alert level (e.g. "info", "warning", "critical").
	Severity string
}

// Sender is the interface implemented by all notification backends.
type Sender interface {
	Send(msg Message) error
}
