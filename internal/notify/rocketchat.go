package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RocketChatClient sends alert messages to a Rocket.Chat incoming webhook.
type RocketChatClient struct {
	webhookURL string
	httpClient *http.Client
}

type rocketChatPayload struct {
	Text string `json:"text"`
}

// NewRocketChatClient constructs a RocketChatClient.
// webhookURL must be a non-empty Rocket.Chat incoming webhook URL.
func NewRocketChatClient(webhookURL string) (*RocketChatClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("rocketchat: webhook URL must not be empty")
	}
	return &RocketChatClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Rocket.Chat webhook.
func (c *RocketChatClient) Send(message string) error {
	payload := rocketChatPayload{Text: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
