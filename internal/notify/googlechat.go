package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GoogleChatClient sends alerts to a Google Chat webhook.
type GoogleChatClient struct {
	webhookURL string
	httpClient *http.Client
}

// NewGoogleChatClient creates a new GoogleChatClient.
// Returns an error if webhookURL is empty.
func NewGoogleChatClient(webhookURL string) (*GoogleChatClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googlechat: webhook URL must not be empty")
	}
	return &GoogleChatClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers msg to the configured Google Chat webhook.
func (c *GoogleChatClient) Send(msg string) error {
	payload := map[string]string{"text": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
