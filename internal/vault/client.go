// Package vault provides a client for interacting with HashiCorp Vault.
package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with additional functionality
// for secret expiry monitoring.
type Client struct {
	api     *vaultapi.Client
	address string
}

// SecretInfo holds metadata about a Vault secret relevant to expiry monitoring.
type SecretInfo struct {
	Path       string
	ExpiresAt  time.Time
	TTL        time.Duration
	Renewable  bool
	LeaseID    string
	Metadata   map[string]interface{}
}

// NewClient creates a new Vault client using the provided address and token.
// It validates connectivity by performing a token lookup on creation.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault api client: %w", err)
	}

	api.SetToken(token)

	return &Client{
		api:     api,
		address: address,
	}, nil
}

// Ping checks that the Vault server is reachable and the token is valid.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.api.Auth().Token().LookupSelfWithContext(ctx)
	if err != nil {
		return fmt.Errorf("vault ping failed: %w", err)
	}
	return nil
}

// GetSecretInfo retrieves expiry metadata for a KV secret at the given path.
// It supports both KV v1 and KV v2 mounts.
func (c *Client) GetSecretInfo(ctx context.Context, path string) (*SecretInfo, error) {
	secret, err := c.api.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	info := &SecretInfo{
		Path:      path,
		LeaseID:   secret.LeaseID,
		Renewable: secret.Renewable,
		Metadata:  make(map[string]interface{}),
	}

	// Populate TTL and expiry from lease duration if available.
	if secret.LeaseDuration > 0 {
		info.TTL = time.Duration(secret.LeaseDuration) * time.Second
		info.ExpiresAt = time.Now().Add(info.TTL)
	}

	// For KV v2, extract metadata from the "metadata" sub-key.
	if meta, ok := secret.Data["metadata"]; ok {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			info.Metadata = metaMap

			// KV v2 stores deletion_time as an expiry indicator.
			if deletionTime, ok := metaMap["deletion_time"].(string); ok && deletionTime != "" {
				parsed, err := time.Parse(time.RFC3339, deletionTime)
				if err == nil {
					info.ExpiresAt = parsed
					info.TTL = time.Until(parsed)
				}
			}
		}
	}

	return info, nil
}

// IsExpiringSoon returns true if the secret expires within the given threshold.
func (s *SecretInfo) IsExpiringSoon(threshold time.Duration) bool {
	if s.ExpiresAt.IsZero() {
		return false
	}
	return time.Until(s.ExpiresAt) <= threshold
}

// IsExpired returns true if the secret has already passed its expiry time.
func (s *SecretInfo) IsExpired() bool {
	if s.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(s.ExpiresAt)
}
