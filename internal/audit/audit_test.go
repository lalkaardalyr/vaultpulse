package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/audit"
)

func TestLog_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.EventScanStarted, "", "scan initiated", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if entry.Event != audit.EventScanStarted {
		t.Errorf("expected event %q, got %q", audit.EventScanStarted, entry.Event)
	}
}

func TestLog_IncludesPath(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Log(audit.EventSecretExpiry, "secret/db/password", "expires soon", nil)

	var entry audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)
	if entry.Path != "secret/db/password" {
		t.Errorf("expected path %q, got %q", "secret/db/password", entry.Path)
	}
}

func TestLog_IncludesMeta(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	meta := map[string]string{"severity": "critical", "days": "2"}
	_ = l.Log(audit.EventAlertSent, "secret/api/key", "alert dispatched", meta)

	var entry audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)
	if entry.Meta["severity"] != "critical" {
		t.Errorf("expected meta severity=critical, got %q", entry.Meta["severity"])
	}
}

func TestLog_EndsWithNewline(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Log(audit.EventScanFinished, "", "done", nil)

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected output to end with newline")
	}
}

func TestNewWithOptions_WithWriter(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithOptions(audit.WithWriter(&buf))

	_ = l.Log(audit.EventScanStarted, "", "test", nil)

	if buf.Len() == 0 {
		t.Error("expected output in provided writer, got nothing")
	}
}
