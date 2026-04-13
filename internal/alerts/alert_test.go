package alerts

import (
	"testing"
	"time"
)

func TestNewAlert_Critical(t *testing.T) {
	warn := 48 * time.Hour
	critical := 24 * time.Hour
	expiry := time.Now().Add(12 * time.Hour)

	a := NewAlert("secret/db/password", expiry, warn, critical)

	if a.Level != LevelCritical {
		t.Errorf("expected CRITICAL, got %s", a.Level)
	}
	if a.SecretPath != "secret/db/password" {
		t.Errorf("unexpected secret path: %s", a.SecretPath)
	}
}

func TestNewAlert_Warning(t *testing.T) {
	warn := 48 * time.Hour
	critical := 24 * time.Hour
	expiry := time.Now().Add(36 * time.Hour)

	a := NewAlert("secret/api/key", expiry, warn, critical)

	if a.Level != LevelWarning {
		t.Errorf("expected WARNING, got %s", a.Level)
	}
}

func TestNewAlert_Info(t *testing.T) {
	warn := 48 * time.Hour
	critical := 24 * time.Hour
	expiry := time.Now().Add(72 * time.Hour)

	a := NewAlert("secret/safe", expiry, warn, critical)

	if a.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", a.Level)
	}
}

func TestNewAlert_MessageNotEmpty(t *testing.T) {
	a := NewAlert("secret/x", time.Now().Add(time.Hour), 48*time.Hour, 24*time.Hour)
	if a.Message == "" {
		t.Error("expected non-empty message")
	}
}
