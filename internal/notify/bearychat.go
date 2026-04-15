package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// BearyChat sends alert messages to a BearyChat incoming webhook.
type BearyChat struct {
	webhookURL string
	client     *http.Client
}

type bearyChatPayload struct {
	Text        string `json:"text"`
	Notification string `json:"notification"`
}

// NewBearyChatClient constructs a BearyChat sender.
// Returns an error if webhookURL is empty.
func NewBearyChatClient(webhookURL string) (*BearyChat, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("bearychat: webhook URL must not be empty")
	}
	return &BearyChat{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send posts the message to the BearyChat webhook endpoint.
func (b *BearyChat) Send(message string) error {
	payload := bearyChatPayload{
		Text:        message,
		Notification: message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: failed to marshal payload: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
