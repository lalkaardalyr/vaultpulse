package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordClient sends alert messages to a Discord channel via webhook.
type DiscordClient struct {
	webhookURL string
	httpClient *http.Client
}

type discordPayload struct {
	Content string `json:"content"`
}

// NewDiscordClient creates a new DiscordClient.
// Returns an error if webhookURL is empty.
func NewDiscordClient(webhookURL string) (*DiscordClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("discord: webhook URL must not be empty")
	}
	return &DiscordClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Discord webhook.
func (d *DiscordClient) Send(message string) error {
	payload := discordPayload{Content: message}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: failed to marshal payload: %w", err)
	}

	resp, err := d.httpClient.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Discord returns 204 No Content on success.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
