package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookClient sends alert notifications to a generic HTTP webhook endpoint.
type WebhookClient struct {
	url        string
	httpClient *http.Client
}

type webhookPayload struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewWebhookClient creates a new WebhookClient.
// Returns an error if the URL is empty.
func NewWebhookClient(url string) (*WebhookClient, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook: url must not be empty")
	}
	return &WebhookClient{
		url: url,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers the alert message to the configured webhook URL.
func (w *WebhookClient) Send(msg Message) error {
	payload := webhookPayload{
		Title:     "VaultPulse Alert",
		Message:   msg.Body,
		Severity:  msg.Severity,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: failed to marshal payload: %w", err)
	}

	resp, err := w.httpClient.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
