package alerts

import (
	"fmt"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo     Level = "INFO"
	LevelWarning  Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Alert holds information about a secret expiry notification.
type Alert struct {
	SecretPath string
	Level      Level
	ExpiresAt  time.Time
	TimeLeft   time.Duration
	Message    string
}

// NewAlert constructs an Alert for the given secret path and expiry time.
func NewAlert(secretPath string, expiresAt time.Time, warnThreshold, criticalThreshold time.Duration) Alert {
	timeLeft := time.Until(expiresAt)
	level := LevelInfo

	switch {
	case timeLeft <= criticalThreshold:
		level = LevelCritical
	case timeLeft <= warnThreshold:
		level = LevelWarning
	}

	msg := fmt.Sprintf("[%s] Secret '%s' expires in %s (at %s)",
		level, secretPath, timeLeft.Round(time.Second), expiresAt.Format(time.RFC3339))

	return Alert{
		SecretPath: secretPath,
		Level:      level,
		ExpiresAt:  expiresAt,
		TimeLeft:   timeLeft,
		Message:    msg,
	}
}
