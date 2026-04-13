package runner_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/example/vaultpulse/internal/config"
	"github.com/example/vaultpulse/internal/runner"
)

func minimalConfig() *config.Config {
	return &config.Config{
		Vault: config.VaultConfig{
			Address: "http://127.0.0.1:8200",
			Token:   "root",
		},
		Output: config.OutputConfig{
			Format: "table",
		},
		Alerts: config.AlertsConfig{
			WarningDays:  7,
			CriticalDays: 2,
		},
	}
}

func TestNew_ReturnsRunner(t *testing.T) {
	cfg := minimalConfig()
	var buf bytes.Buffer

	r, err := runner.New(cfg, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestNew_InvalidAddress_ReturnsError(t *testing.T) {
	cfg := minimalConfig()
	cfg.Vault.Address = "://bad address"
	var buf bytes.Buffer

	_, err := runner.New(cfg, &buf)
	if err == nil {
		t.Fatal("expected error for invalid vault address, got nil")
	}
}

func TestRun_TableFormat_NoError(t *testing.T) {
	cfg := minimalConfig()
	cfg.Vault.Paths = []string{} // empty path list — monitor returns nothing
	var buf bytes.Buffer

	r, err := runner.New(cfg, &buf)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	ctx := context.Background()
	if err := r.Run(ctx); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
}

func TestRun_JSONFormat_NoError(t *testing.T) {
	cfg := minimalConfig()
	cfg.Output.Format = "json"
	cfg.Vault.Paths = []string{}
	var buf bytes.Buffer

	r, err := runner.New(cfg, &buf)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("Run (json) returned unexpected error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty JSON output")
	}
}
