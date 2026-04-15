package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MattermostClient sends alert messages to a Mattermost incoming webhook.
type MattermostClient struct {
	webhookURL string
	httpClient *http.Client
}

type mattermostPayload struct {
	Text     string `json:"text"`
	Username string `json:"username,omitempty"`
	Channel  string `json:"channel,omitempty"`
}

// NewMattermostClient creates a new MattermostClient.
// webhookURL must be a valid Mattermost incoming webhook URL.
func NewMattermostClient(webhookURL string) (*MattermostClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("mattermost: webhook URL must not be empty")
	}
	return &MattermostClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Mattermost webhook.
func (c *MattermostClient) Send(message string) error {
	payload := mattermostPayload{
		Text:     message,
		Username: "VaultPulse",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
