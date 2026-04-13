// Package audit provides structured audit logging for VaultPulse
// secret monitoring events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType classifies the kind of audit event.
type EventType string

const (
	EventScanStarted  EventType = "scan_started"
	EventScanFinished EventType = "scan_finished"
	EventAlertSent    EventType = "alert_sent"
	EventSecretExpiry EventType = "secret_expiry"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time  `json:"timestamp"`
	Event     EventType  `json:"event"`
	Path      string     `json:"path,omitempty"`
	Message   string     `json:"message"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes audit entries to an io.Writer as newline-delimited JSON.
type Logger struct {
	writer io.Writer
}

// New returns a Logger that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{writer: w}
}

// Log writes a single audit entry.
func (l *Logger) Log(event EventType, path, message string, meta map[string]string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Path:      path,
		Message:   message,
		Meta:      meta,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.writer, "%s\n", b)
	return err
}
