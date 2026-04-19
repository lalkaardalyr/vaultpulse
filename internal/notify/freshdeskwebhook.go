package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FreshdeskWebhookClient sends alert notifications to a generic Freshdesk-style inbound webhook.
type FreshdeskWebhookClient struct {
	webhookURL string
	apiKey     string
	httpClient *http.Client
}

type freshdeskWebhookPayload struct {
	Text string `json:"text"`
}

// NewFreshdeskWebhookClient creates a new FreshdeskWebhookClient.
func NewFreshdeskWebhookClient(webhookURL, apiKey string) (*FreshdeskWebhookClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("freshdeskwebhook: webhook URL must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("freshdeskwebhook: API key must not be empty")
	}
	return &FreshdeskWebhookClient{
		webhookURL: webhookURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert message to the Freshdesk webhook endpoint.
func (c *FreshdeskWebhookClient) Send(message string) error {
	payload := freshdeskWebhookPayload{Text: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("freshdeskwebhook: marshal error: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("freshdeskwebhook: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("freshdeskwebhook: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshdeskwebhook: unexpected status: %d", resp.StatusCode)
	}
	return nil
}
