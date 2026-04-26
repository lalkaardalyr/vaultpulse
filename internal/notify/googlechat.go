package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GoogleChatClient sends alert messages to a Google Chat webhook.
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

// Send posts the message to the Google Chat webhook.
func (c *GoogleChatClient) Send(message string) error {
	payload := map[string]string{"text": message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
