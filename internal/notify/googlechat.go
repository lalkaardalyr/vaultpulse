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

// NewGoogleChatClient constructs a GoogleChatClient.
func NewGoogleChatClient(webhookURL string) (*GoogleChatClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googlechat: webhook URL must not be empty")
	}
	return &GoogleChatClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts a message to the Google Chat webhook.
func (c *GoogleChatClient) Send(msg string) error {
	payload := map[string]string{"text": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal error: %w", err)
	}
	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
