package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// JumpCloudClient sends alert events to JumpCloud via its webhook API.
type JumpCloudClient struct {
	webhookURL string
	apiKey     string
	httpClient *http.Client
}

// NewJumpCloudClient creates a new JumpCloudClient.
// Both webhookURL and apiKey must be non-empty.
func NewJumpCloudClient(webhookURL, apiKey string) (*JumpCloudClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("jumpcloud: webhook URL must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("jumpcloud: API key must not be empty")
	}
	return &JumpCloudClient{
		webhookURL: webhookURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert message to the JumpCloud webhook endpoint.
func (c *JumpCloudClient) Send(message string) error {
	payload := map[string]string{
		"message": message,
		"source":  "vaultpulse",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("jumpcloud: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jumpcloud: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("jumpcloud: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jumpcloud: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
