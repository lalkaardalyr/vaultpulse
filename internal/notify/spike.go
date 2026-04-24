package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SpikeClient sends alerts to Spike.sh via an incoming webhook.
type SpikeClient struct {
	webhookURL string
	httpClient *http.Client
}

// NewSpikeClient creates a new SpikeClient.
// Returns an error if webhookURL is empty.
func NewSpikeClient(webhookURL string) (*SpikeClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("spike: webhook URL must not be empty")
	}
	return &SpikeClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers msg to the configured Spike.sh webhook.
func (c *SpikeClient) Send(msg string) error {
	payload := map[string]string{"message": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("spike: marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("spike: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("spike: unexpected status %d", resp.StatusCode)
	}
	return nil
}
