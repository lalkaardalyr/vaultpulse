package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const clickatellDefaultEndpoint = "https://platform.clickatell.com/messages"

type clickatellClient struct {
	apiKey   string
	to       string
	endpoint string
	httpClient *http.Client
}

type clickatellPayload struct {
	Messages []clickatellMessage `json:"messages"`
}

type clickatellMessage struct {
	Channel string `json:"channel"`
	To      string `json:"to"`
	Content string `json:"content"`
}

// NewClickatellClient creates a new Clickatell SMS notification client.
func NewClickatellClient(apiKey, to string) (Sender, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("clickatell: api key must not be empty")
	}
	if to == "" {
		return nil, fmt.Errorf("clickatell: recipient phone number must not be empty")
	}
	return &clickatellClient{
		apiKey:     apiKey,
		to:         to,
		endpoint:   clickatellDefaultEndpoint,
		httpClient: &http.Client{},
	}, nil
}

func (c *clickatellClient) Send(ctx context.Context, message string) error {
	payload := clickatellPayload{
		Messages: []clickatellMessage{
			{Channel: "sms", To: c.to, Content: message},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("clickatell: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("clickatell: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("clickatell: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clickatell: unexpected status %d", resp.StatusCode)
	}
	return nil
}
