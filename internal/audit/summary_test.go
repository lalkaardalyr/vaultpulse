package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpulse/internal/audit"
)

func TestSummaryLogger_Log_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewSummaryLogger(&buf)

	err := logger.Log(10, 2, 3, 5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.SummaryEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestSummaryLogger_Log_CountsAreCorrect(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewSummaryLogger(&buf)

	_ = logger.Log(9, 1, 4, 4, nil)

	var entry audit.SummaryEntry
	_ = json.Unmarshal(buf.Bytes(), &entry)

	if entry.Total != 9 {
		t.Errorf("expected total=9, got %d", entry.Total)
	}
	if entry.Critical != 1 {
		t.Errorf("expected critical=1, got %d", entry.Critical)
	}
	if entry.Warning != 4 {
		t.Errorf("expected warning=4, got %d", entry.Warning)
	}
	if entry.Healthy != 4 {
		t.Errorf("expected healthy=4, got %d", entry.Healthy)
	}
}

func TestSummaryLogger_Log_IncludesMeta(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewSummaryLogger(&buf)

	meta := map[string]any{"run_id": "abc-123"}
	_ = logger.Log(1, 0, 0, 1, meta)

	if !strings.Contains(buf.String(), "run_id") {
		t.Error("expected meta key 'run_id' in output")
	}
}

func TestSummaryLogger_Log_EndsWithNewline(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewSummaryLogger(&buf)

	_ = logger.Log(0, 0, 0, 0, nil)

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected output to end with newline")
	}
}

func TestSummaryLogger_Log_TimestampPresent(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewSummaryLogger(&buf)

	_ = logger.Log(3, 1, 1, 1, nil)

	var entry audit.SummaryEntry
	_ = json.Unmarshal(buf.Bytes(), &entry)

	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp in summary entry")
	}
}
