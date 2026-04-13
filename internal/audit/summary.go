package audit

import (
	"encoding/json"
	"io"
	"time"
)

// SummaryEntry represents an aggregated summary of a monitoring run.
type SummaryEntry struct {
	Timestamp  time.Time      `json:"timestamp"`
	Total      int            `json:"total"`
	Critical   int            `json:"critical"`
	Warning    int            `json:"warning"`
	Healthy    int            `json:"healthy"`
	Meta       map[string]any `json:"meta,omitempty"`
}

// SummaryLogger writes structured summary entries to an io.Writer.
type SummaryLogger struct {
	w io.Writer
}

// NewSummaryLogger creates a new SummaryLogger that writes to w.
func NewSummaryLogger(w io.Writer) *SummaryLogger {
	return &SummaryLogger{w: w}
}

// Log writes a JSON-encoded SummaryEntry to the underlying writer.
// It appends a newline after each entry.
func (s *SummaryLogger) Log(total, critical, warning, healthy int, meta map[string]any) error {
	entry := SummaryEntry{
		Timestamp: time.Now().UTC(),
		Total:     total,
		Critical:  critical,
		Warning:   warning,
		Healthy:   healthy,
		Meta:      meta,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = s.w.Write(data)
	return err
}
