package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SIGNL4Client sends alerts to SIGNL4 via webhook.
type SIGNL4Client struct {
	webhookURL string
	httpClient *http.Client
}

type signl4Payload struct {
	Title   string `json:"Title"`
	Message string `json:"Message"`
	Severity int    `json:"X-S4-SourceSystem,omitempty"`
}

// NewSIGNL4Client creates a new SIGNL4Client.
// webhookURL must be the full SIGNL4 team secret URL.
func NewSIGNL4Client(webhookURL string) (*SIGNL4Client, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("signl4: webhook URL must not be empty")
	}
	return &SIGNL4Client{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the alert message to SIGNL4.
func (c *SIGNL4Client) Send(message string) error {
	payload := signl4Payload{
		Title:   "VaultPulse Alert",
		Message: message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signl4: failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signl4: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signl4: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
