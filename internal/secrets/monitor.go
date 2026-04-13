package secrets

import (
	"fmt"
	"time"

	"github.com/vaultpulse/internal/vault"
)

// SecretStatus represents the expiry status of a monitored secret.
type SecretStatus struct {
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
	Status    string // "ok", "warning", "critical", "expired"
}

// Monitor checks secrets in Vault and returns their expiry statuses.
type Monitor struct {
	client          *vault.Client
	warningThreshold time.Duration
	criticalThreshold time.Duration
}

// NewMonitor creates a new Monitor with the given Vault client and thresholds.
func NewMonitor(client *vault.Client, warningThreshold, criticalThreshold time.Duration) *Monitor {
	return &Monitor{
		client:            client,
		warningThreshold:  warningThreshold,
		criticalThreshold: criticalThreshold,
	}
}

// CheckSecret retrieves a secret's lease info and determines its status.
func (m *Monitor) CheckSecret(path string) (*SecretStatus, error) {
	ttl, err := m.client.GetSecretTTL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get TTL for secret %q: %w", path, err)
	}

	now := time.Now()
	expiresAt := now.Add(ttl)

	status := classifyStatus(ttl, m.warningThreshold, m.criticalThreshold)

	return &SecretStatus{
		Path:      path,
		ExpiresAt: expiresAt,
		TTL:       ttl,
		Status:    status,
	}, nil
}

// CheckAll checks all provided secret paths and returns a slice of statuses.
func (m *Monitor) CheckAll(paths []string) ([]*SecretStatus, error) {
	statuses := make([]*SecretStatus, 0, len(paths))
	for _, path := range paths {
		s, err := m.CheckSecret(path)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

func classifyStatus(ttl, warning, critical time.Duration) string {
	switch {
	case ttl <= 0:
		return "expired"
	case ttl <= critical:
		return "critical"
	case ttl <= warning:
		return "warning"
	default:
		return "ok"
	}
}
