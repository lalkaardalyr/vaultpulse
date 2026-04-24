package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RevoltClient sends alert messages to a Revolt webhook.
type RevoltClient struct {
	webhookURL string
	httpClient *http.Client
}

type revoltPayload struct {
	Content string `json:"content"`
}

// NewRevoltClient creates a new RevoltClient.
// Returns an error if webhookURL is empty.
func NewRevoltClient(webhookURL string) (*RevoltClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("revolt: webhook URL must not be empty")
	}
	return &RevoltClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Revolt webhook.
func (c *RevoltClient) Send(message string) error {
	payload := revoltPayload{Content: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("revolt: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("revolt: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("revolt: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
