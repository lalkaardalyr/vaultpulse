package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CustomEventClient sends alerts to a generic custom event endpoint
// using a configurable HTTP POST with a JSON payload.
type CustomEventClient struct {
	url     string
	apiKey  string
	source  string
	httpClient *http.Client
}

type customEventPayload struct {
	Source  string `json:"source"`
	Message string `json:"message"`
}

// NewCustomEventClient creates a new CustomEventClient.
func NewCustomEventClient(url, apiKey, source string) (*CustomEventClient, error) {
	if url == "" {
		return nil, fmt.Errorf("customevent: url must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("customevent: api key must not be empty")
	}
	if source == "" {
		return nil, fmt.Errorf("customevent: source must not be empty")
	}
	return &CustomEventClient{
		url:        url,
		apiKey:     apiKey,
		source:     source,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the alert message to the custom event endpoint.
func (c *CustomEventClient) Send(message string) error {
	payload := customEventPayload{
		Source:  c.source,
		Message: message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("customevent: failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("customevent: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("customevent: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("customevent: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
