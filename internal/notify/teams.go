package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TeamsClient sends alert messages to a Microsoft Teams channel via
// an incoming webhook URL.
type TeamsClient struct {
	webhookURL string
	httpClient *http.Client
}

type teamsPayload struct {
	Text string `json:"text"`
}

// NewTeamsClient constructs a TeamsClient for the given webhook URL.
// It returns an error if the URL is empty.
func NewTeamsClient(webhookURL string) (*TeamsClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("teams: webhook URL must not be empty")
	}
	return &TeamsClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Teams webhook.
func (c *TeamsClient) Send(message string) error {
	payload := teamsPayload{Text: message}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
