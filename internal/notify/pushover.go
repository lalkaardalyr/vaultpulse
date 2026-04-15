package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverClient sends alerts via the Pushover notification service.
type PushoverClient struct {
	apiToken  string
	userKey   string
	apiURL    string
	httpClient *http.Client
}

type pushoverPayload struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
	Title   string `json:"title,omitempty"`
}

// NewPushoverClient creates a new PushoverClient.
// apiToken is the application API token and userKey is the recipient user key.
func NewPushoverClient(apiToken, userKey string) (*PushoverClient, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("pushover: api token must not be empty")
	}
	if userKey == "" {
		return nil, fmt.Errorf("pushover: user key must not be empty")
	}
	return &PushoverClient{
		apiToken:   apiToken,
		userKey:    userKey,
		apiURL:     pushoverAPIURL,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the message via Pushover.
func (c *PushoverClient) Send(message string) error {
	payload := pushoverPayload{
		Token:   c.apiToken,
		User:    c.userKey,
		Message: message,
		Title:   "VaultPulse Alert",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushover: failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushover: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushover: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
