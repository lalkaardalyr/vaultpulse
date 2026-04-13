// Package notify provides integrations for sending alerts to external
// notification services such as Slack.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackClient sends alert messages to a Slack webhook URL.
type SlackClient struct {
	webhookURL string
	httpClient *http.Client
}

// slackPayload is the JSON body sent to the Slack Incoming Webhook API.
type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackClient creates a SlackClient that posts to the given webhook URL.
// A default HTTP timeout of 10 seconds is applied.
func NewSlackClient(webhookURL string) (*SlackClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("notify: slack webhook URL must not be empty")
	}
	return &SlackClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts the given message text to the configured Slack webhook.
// It returns an error if the HTTP request fails or the server responds
// with a non-2xx status code.
func (s *SlackClient) Send(message string) error {
	payload := slackPayload{Text: message}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: failed to marshal slack payload: %w", err)
	}

	resp, err := s.httpClient.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: slack HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: slack responded with status %d", resp.StatusCode)
	}

	return nil
}
