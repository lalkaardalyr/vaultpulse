package secrets

import (
	"testing"
	"time"
)

func TestClassifyStatus(t *testing.T) {
	warning := 48 * time.Hour
	critical := 24 * time.Hour

	tests := []struct {
		name     string
		ttl      time.Duration
		expected string
	}{
		{"expired", -1 * time.Second, "expired"},
		{"zero ttl", 0, "expired"},
		{"critical", 12 * time.Hour, "critical"},
		{"critical boundary", 24 * time.Hour, "critical"},
		{"warning", 36 * time.Hour, "warning"},
		{"warning boundary", 48 * time.Hour, "warning"},
		{"ok", 72 * time.Hour, "ok"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyStatus(tc.ttl, warning, critical)
			if got != tc.expected {
				t.Errorf("classifyStatus(%v) = %q, want %q", tc.ttl, got, tc.expected)
			}
		})
	}
}

func TestSecretStatus_Fields(t *testing.T) {
	now := time.Now()
	ttl := 72 * time.Hour

	s := &SecretStatus{
		Path:      "secret/data/myapp/db",
		ExpiresAt: now.Add(ttl),
		TTL:       ttl,
		Status:    "ok",
	}

	if s.Path != "secret/data/myapp/db" {
		t.Errorf("unexpected Path: %s", s.Path)
	}
	if s.Status != "ok" {
		t.Errorf("unexpected Status: %s", s.Status)
	}
	if s.TTL != ttl {
		t.Errorf("unexpected TTL: %v", s.TTL)
	}
	if s.ExpiresAt.Before(now) {
		t.Errorf("ExpiresAt should be in the future")
	}
}

func TestNewMonitor(t *testing.T) {
	m := NewMonitor(nil, 48*time.Hour, 24*time.Hour)
	if m == nil {
		t.Fatal("expected non-nil Monitor")
	}
	if m.warningThreshold != 48*time.Hour {
		t.Errorf("unexpected warningThreshold: %v", m.warningThreshold)
	}
	if m.criticalThreshold != 24*time.Hour {
		t.Errorf("unexpected criticalThreshold: %v", m.criticalThreshold)
	}
}
