package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultpulse/internal/audit"
)

func TestRotationLogger_Log_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	rl := audit.NewRotationLogger(&buf)

	if err := rl.Log("secret/db", "critical", true, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var event audit.RotationEvent
	if err := json.Unmarshal(buf.Bytes(), &event); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestRotationLogger_Log_ContainsPath(t *testing.T) {
	var buf bytes.Buffer
	rl := audit.NewRotationLogger(&buf)

	_ = rl.Log("secret/api-key", "warning", false, nil)

	if !strings.Contains(buf.String(), "secret/api-key") {
		t.Errorf("expected path in output, got: %s", buf.String())
	}
}

func TestRotationLogger_Log_TriggeredField(t *testing.T) {
	var buf bytes.Buffer
	rl := audit.NewRotationLogger(&buf)

	_ = rl.Log("secret/token", "critical", true, nil)

	var event audit.RotationEvent
	_ = json.Unmarshal(buf.Bytes(), &event)

	if !event.Triggered {
		t.Errorf("expected triggered=true, got false")
	}
}

func TestRotationLogger_Log_MetaIncluded(t *testing.T) {
	var buf bytes.Buffer
	rl := audit.NewRotationLogger(&buf)

	meta := map[string]string{"owner": "platform-team"}
	_ = rl.Log("secret/cert", "info", false, meta)

	var event audit.RotationEvent
	_ = json.Unmarshal(buf.Bytes(), &event)

	if event.Meta["owner"] != "platform-team" {
		t.Errorf("expected meta owner=platform-team, got %v", event.Meta)
	}
}

func TestRotationLogger_Log_EndsWithNewline(t *testing.T) {
	var buf bytes.Buffer
	rl := audit.NewRotationLogger(&buf)

	_ = rl.Log("secret/x", "info", false, nil)

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("expected output to end with newline")
	}
}

func TestNewRotationLogger_NilWriter_UsesStdout(t *testing.T) {
	rl := audit.NewRotationLogger(nil)
	if rl == nil {
		t.Fatal("expected non-nil RotationLogger")
	}
}
