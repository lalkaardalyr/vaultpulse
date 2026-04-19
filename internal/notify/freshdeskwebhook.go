package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FreshdeskWebhookClient sends alerts via a Freshdesk webhook URL.
type FreshdeskWebhookClient struct {
	webhookURL string
	httpClient *http.Client
}

// NewFreshdeskWebhookClient creates a new FreshdeskWebhookClient.
func NewFreshdeskWebhookClient(webhookURL string) (*FreshdeskWebhookClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("freshdeskwebhook: webhook URL must not be empty")
	}
	return &FreshdeskWebhookClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts a message payload to the Freshdesk webhook.
func (c *FreshdeskWebhookClient) Send(msg string) error {
	payload := map[string]string{"text": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("freshdeskwebhook: marshal error: %w", err)
	}
	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("freshdeskwebhook: request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshdeskwebhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
