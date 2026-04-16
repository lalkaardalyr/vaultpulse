package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const clickSendBaseURL = "https://rest.clicksend.com/v3/sms/send"

// ClickSendClient sends SMS alerts via the ClickSend API.
type ClickSendClient struct {
	username string
	apiKey   string
	to       string
	httpClient *http.Client
}

// NewClickSendClient creates a new ClickSendClient.
func NewClickSendClient(username, apiKey, to string) (*ClickSendClient, error) {
	if username == "" {
		return nil, fmt.Errorf("clicksend: username must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("clicksend: api key must not be empty")
	}
	if to == "" {
		return nil, fmt.Errorf("clicksend: recipient phone number must not be empty")
	}
	return &ClickSendClient{
		username:   username,
		apiKey:     apiKey,
		to:         to,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the message via ClickSend SMS API.
func (c *ClickSendClient) Send(message string) error {
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"to": c.to, "body": message, "source": "vaultpulse"},
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("clicksend: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, clickSendBaseURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("clicksend: build request: %w", err)
	}
	req.SetBasicAuth(c.username, c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("clicksend: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clicksend: unexpected status %d", resp.StatusCode)
	}
	return nil
}
