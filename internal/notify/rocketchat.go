package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RocketChatClient sends alerts to a Rocket.Chat webhook.
type RocketChatClient struct {
	webhookURL string
	httpClient *http.Client
}

// NewRocketChatClient creates a new RocketChatClient.
func NewRocketChatClient(webhookURL string) (*RocketChatClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("rocketchat: webhook URL must not be empty")
	}
	return &RocketChatClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts a message to Rocket.Chat.
func (c *RocketChatClient) Send(msg string) error {
	payload := map[string]string{"text": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: marshal error: %w", err)
	}
	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
