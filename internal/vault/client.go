package vault

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	vc      *vaultapi.Client
	address string
	token   string
}

// NewClient creates and configures a new Vault client.
func NewClient(address, token string, timeout time.Duration) (*Client, error) {
	if address == "" {
		return nil, errors.New("vault address must not be empty")
	}
	if token == "" {
		return nil, errors.New("vault token must not be empty")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address
	cfg.HttpClient = &http.Client{Timeout: timeout}

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault api client: %w", err)
	}
	vc.SetToken(token)

	return &Client{
		vc:      vc,
		address: address,
		token:   token,
	}, nil
}

// ReadSecret reads a secret at the given path.
func (c *Client) ReadSecret(path string) (*vaultapi.Secret, error) {
	secret, err := c.vc.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}
	return secret, nil
}

// GetSecretTTL returns the remaining TTL for a secret at the given path.
func (c *Client) GetSecretTTL(path string) (time.Duration, error) {
	secret, err := c.ReadSecret(path)
	if err != nil {
		return 0, err
	}

	if secret.LeaseDuration == 0 {
		// Check KV v2 metadata for expiration
		if meta, ok := secret.Data["metadata"]; ok {
			if metaMap, ok := meta.(map[string]interface{}); ok {
				if expRaw, ok := metaMap["deletion_time"]; ok {
					if expStr, ok := expRaw.(string); ok && expStr != "" {
						exp, err := time.Parse(time.RFC3339, expStr)
						if err == nil {
							return time.Until(exp), nil
						}
					}
				}
			}
		}
		return 0, fmt.Errorf("secret at %q has no lease or expiry information", path)
	}

	return time.Duration(secret.LeaseDuration) * time.Second, nil
}

// Address returns the configured Vault server address.
func (c *Client) Address() string {
	return c.address
}
