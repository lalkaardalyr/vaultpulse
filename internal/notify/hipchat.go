package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HipChatClient sends alert notifications to a HipChat room.
type HipChatClient struct {
	roomURL string
	token   string
	hc      *http.Client
}

type hipChatPayload struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
}

// NewHipChatClient constructs a HipChatClient.
// roomURL is the full HipChat v2 room notification URL.
func NewHipChatClient(roomURL, token string) (*HipChatClient, error) {
	if roomURL == "" {
		return nil, fmt.Errorf("hipchat: room URL must not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("hipchat: token must not be empty")
	}
	return &HipChatClient{
		roomURL: roomURL,
		token:   token,
		hc:      &http.Client{},
	}, nil
}

// Send posts an alert message to the configured HipChat room.
func (c *HipChatClient) Send(msg string) error {
	payload := hipChatPayload{
		Message:       msg,
		MessageFormat: "text",
		Color:         "red",
		Notify:        true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("hipchat: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.roomURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("hipchat: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("hipchat: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
