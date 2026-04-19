package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SpikeClient sends alerts to Spike.sh via webhook.
type SpikeClient struct {
	webhookURL string
	httpClient *http.Client
}

// NewSpikeClient creates a new SpikeClient.
func NewSpikeClient(webhookURL string) (*SpikeClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("spike: webhook URL must not be empty")
	}
	return &SpikeClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert message to Spike.sh.
func (c *SpikeClient) Send(msg string) error {
	payload := map[string]string{"message": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("spike: failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("spike: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("spike: unexpected status %d", resp.StatusCode)
	}
	return nil
}
