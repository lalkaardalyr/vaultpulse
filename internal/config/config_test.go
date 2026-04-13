package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "vaultpulse-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
monitor:
  paths:
    - "secret/data/myapp"
  interval: 30m
alerts:
  warning_days: 10
  critical_days: 2
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault address: %s", cfg.Vault.Address)
	}
	if cfg.Alerts.WarningDays != 10 {
		t.Errorf("expected warning_days=10, got %d", cfg.Alerts.WarningDays)
	}
	if cfg.Monitor.Interval != 30*time.Minute {
		t.Errorf("expected interval=30m, got %v", cfg.Monitor.Interval)
	}
}

func TestLoad_AppliesDefaults(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
monitor:
  paths:
    - "secret/data/myapp"
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Alerts.WarningDays != 14 {
		t.Errorf("expected default warning_days=14, got %d", cfg.Alerts.WarningDays)
	}
	if cfg.Alerts.CriticalDays != 3 {
		t.Errorf("expected default critical_days=3, got %d", cfg.Alerts.CriticalDays)
	}
	if cfg.Monitor.Interval != time.Hour {
		t.Errorf("expected default interval=1h, got %v", cfg.Monitor.Interval)
	}
}

func TestLoad_MissingVaultAddress(t *testing.T) {
	content := `
vault:
  token: "root"
monitor:
  paths:
    - "secret/data/myapp"
`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestLoad_TokenFromEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "env-token")
	content := `
vault:
  address: "http://127.0.0.1:8200"
monitor:
  paths:
    - "secret/data/myapp"
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Vault.Token != "env-token" {
		t.Errorf("expected token from env, got %s", cfg.Vault.Token)
	}
}
