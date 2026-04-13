package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/vaultpulse/internal/output"
	"github.com/vaultpulse/internal/secrets"
)

func sampleStatuses() []secrets.SecretStatus {
	now := time.Now()
	return []secrets.SecretStatus{
		{Path: "secret/db/prod", Status: "critical", ExpiresAt: now.Add(12 * time.Hour)},
		{Path: "secret/api/key", Status: "warning", ExpiresAt: now.Add(5 * 24 * time.Hour)},
		{Path: "secret/tls/cert", Status: "ok", ExpiresAt: now.Add(30 * 24 * time.Hour)},
	}
}

func TestFormatter_WriteTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatTable).WithWriter(&buf)
	if err := f.Write(sampleStatuses()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, header := range []string{"PATH", "STATUS", "EXPIRES IN", "EXPIRY"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in table output", header)
		}
	}
}

func TestFormatter_WriteTable_ContainsPaths(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatTable).WithWriter(&buf)
	_ = f.Write(sampleStatuses())
	out := buf.String()
	if !strings.Contains(out, "secret/db/prod") {
		t.Error("expected secret path in table output")
	}
}

func TestFormatter_WriteJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatJSON).WithWriter(&buf)
	if err := f.Write(sampleStatuses()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Error("JSON output should start with '['") 
	}
	if !strings.Contains(out, "\"path\"") {
		t.Error("JSON output should contain 'path' key")
	}
	if !strings.Contains(out, "\"status\"") {
		t.Error("JSON output should contain 'status' key")
	}
}

func TestFormatter_WriteJSON_EmptySlice(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatJSON).WithWriter(&buf)
	if err := f.Write([]secrets.SecretStatus{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "[\n]" && out != "[]" && !strings.Contains(out, "[") {
		t.Errorf("unexpected empty JSON output: %q", out)
	}
}
