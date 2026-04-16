// Package notify provides notification clients for various platforms.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// LineClient sends notifications via LINE Notify.
type LineClient struct {
	token   string
	httpClient *http.Client
}

// NewLineClient creates a new LINE Notify client.
func NewLineClient(token string) (*LineClient, error) {
	if token == "" {
		return nil, fmt.Errorf("line notify: token must not be empty")
	}
	return &LineClient{
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers a message via LINE Notify API.
func (c *LineClient) Send(msg string) error {
	payload := map[string]string{"message": msg}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("line notify: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://notify-api.line.me/api/notify", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("line notify: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("line notify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("line notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
