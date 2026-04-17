package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SpikeClient sends alert notifications to Spike.sh via their incoming webhook.
type SpikeClient struct {
	webhookURL string
	httpClient *http.Client
}

type spikePayload struct {
	Message string `json:"message"`
	Severity string `json:"severity"`
}

// NewSpikeClient creates a new SpikeClient. webhookURL must not be empty.
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
func (c *SpikeClient) Send(message string) error {
	payload := spikePayload{
		Message:  message,
		Severity: "critical",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("spike: failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("spike: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("spike: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
