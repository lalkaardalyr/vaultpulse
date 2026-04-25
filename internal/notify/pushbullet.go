package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const pushbulletAPIURL = "https://api.pushbullet.com/v2/pushes"

// PushbulletClient sends notifications via the Pushbullet API.
type PushbulletClient struct {
	apiKey string
	url    string
	http   *http.Client
}

type pushbulletPayload struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// NewPushbulletClient creates a new PushbulletClient.
// Returns an error if apiKey is empty.
func NewPushbulletClient(apiKey string) (*PushbulletClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("pushbullet: api key must not be empty")
	}
	return &PushbulletClient{
		apiKey: apiKey,
		url:    pushbulletAPIURL,
		http:   &http.Client{},
	}, nil
}

// Send delivers a notification message via Pushbullet.
func (c *PushbulletClient) Send(message string) error {
	payload := pushbulletPayload{
		Type:  "note",
		Title: "VaultPulse Alert",
		Body:  message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushbullet: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushbullet: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("pushbullet: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushbullet: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
