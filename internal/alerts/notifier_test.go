package alerts

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNotifier_Notify_WritesOnWarning(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(48*time.Hour, 24*time.Hour, WithWriter(&buf))

	expiry := time.Now().Add(36 * time.Hour)
	if err := n.Notify("secret/api", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "WARNING") {
		t.Errorf("expected WARNING in output, got: %s", buf.String())
	}
}

func TestNotifier_Notify_SilentOnInfo(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(48*time.Hour, 24*time.Hour, WithWriter(&buf))

	expiry := time.Now().Add(96 * time.Hour)
	if err := n.Notify("secret/safe", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected no output for INFO level, got: %s", buf.String())
	}
}

func TestNotifier_NotifyAll(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(48*time.Hour, 24*time.Hour, WithWriter(&buf))

	secrets := map[string]time.Time{
		"secret/a": time.Now().Add(10 * time.Hour),
		"secret/b": time.Now().Add(36 * time.Hour),
		"secret/c": time.Now().Add(96 * time.Hour),
	}

	if err := n.NotifyAll(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "CRITICAL") {
		t.Error("expected CRITICAL alert in output")
	}
	if !strings.Contains(output, "WARNING") {
		t.Error("expected WARNING alert in output")
	}
}
