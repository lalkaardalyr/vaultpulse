package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// LarkClient sends alert messages to a Lark (Feishu) incoming webhook.
// Lark webhooks accept a JSON payload with a "msg_type" and "content" field.
type LarkClient struct {
	webhookURL string
	httpClient *http.Client
}

// larkPayload represents the JSON body sent to the Lark webhook endpoint.
type larkPayload struct {
	MsgType string      `json:"msg_type"`
	Content larkContent `json:"content"`
}

type larkContent struct {
	Text string `json:"text"`
}

// NewLarkClient creates a new LarkClient for the given webhook URL.
// Returns an error if the URL is empty.
func NewLarkClient(webhookURL string) (*LarkClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("lark: webhook URL must not be empty")
	}
	return &LarkClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the alert message to the configured Lark webhook.
// The message is wrapped in a plain-text Lark card payload.
func (c *LarkClient) Send(message string) error {
	payload := larkPayload{
		MsgType: "text",
		Content: larkContent{Text: message},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lark: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lark: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
