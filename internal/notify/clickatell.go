package notify

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const clickatellEndpoint = "https://platform.clickatell.com/messages/http/send"

// ClickatellClient sends SMS alerts via the Clickatell HTTP API.
type ClickatellClient struct {
	apiKey string
	to     string
	url    string
}

// NewClickatellClient creates a new ClickatellClient.
func NewClickatellClient(apiKey, to string) (*ClickatellClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("clickatell: api key must not be empty")
	}
	if strings.TrimSpace(to) == "" {
		return nil, fmt.Errorf("clickatell: recipient must not be empty")
	}
	return &ClickatellClient{apiKey: apiKey, to: to, url: clickatellEndpoint}, nil
}

// Send delivers a message via Clickatell.
func (c *ClickatellClient) Send(ctx context.Context, message string) error {
	payload := fmt.Sprintf(`{"apiKey":%q,"to":[%q],"content":%q}`, c.apiKey, c.to, message)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("clickatell: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("clickatell: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clickatell: unexpected status %d", resp.StatusCode)
	}
	return nil
}
