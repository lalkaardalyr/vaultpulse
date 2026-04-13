package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// RotationEvent represents a recorded secret rotation audit event.
type RotationEvent struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Status    string            `json:"status"`
	Triggered bool              `json:"triggered"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// RotationLogger writes rotation audit events as newline-delimited JSON.
type RotationLogger struct {
	w io.Writer
}

// NewRotationLogger returns a RotationLogger writing to w.
// If w is nil, os.Stdout is used.
func NewRotationLogger(w io.Writer) *RotationLogger {
	if w == nil {
		w = os.Stdout
	}
	return &RotationLogger{w: w}
}

// Log writes a single RotationEvent to the underlying writer.
func (rl *RotationLogger) Log(path, status string, triggered bool, meta map[string]string) error {
	event := RotationEvent{
		Timestamp: time.Now().UTC(),
		Path:      path,
		Status:    status,
		Triggered: triggered,
		Meta:      meta,
	}

	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("audit: marshal rotation event: %w", err)
	}

	_, err = fmt.Fprintf(rl.w, "%s\n", b)
	return err
}
