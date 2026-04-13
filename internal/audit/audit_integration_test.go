package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourusername/vaultpulse/internal/audit"
)

// TestLog_MultipleEntries verifies that each call appends a separate JSON line.
func TestLog_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	events := []audit.EventType{
		audit.EventScanStarted,
		audit.EventSecretExpiry,
		audit.EventAlertSent,
		audit.EventScanFinished,
	}
	for _, ev := range events {
		if err := l.Log(ev, "secret/test", "msg", nil); err != nil {
			t.Fatalf("Log(%q) error: %v", ev, err)
		}
	}

	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	if len(lines) != len(events) {
		t.Fatalf("expected %d lines, got %d", len(events), len(lines))
	}

	for i, line := range lines {
		var entry audit.Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
		if entry.Event != events[i] {
			t.Errorf("line %d: expected event %q, got %q", i, events[i], entry.Event)
		}
	}
}

// TestLog_TimestampIsUTC ensures timestamps are recorded in UTC.
func TestLog_TimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.Log(audit.EventScanStarted, "", "tz check", nil)

	var entry audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)

	if entry.Timestamp.Location().String() != "UTC" {
		t.Errorf("expected UTC timestamp, got %s", entry.Timestamp.Location())
	}
}
