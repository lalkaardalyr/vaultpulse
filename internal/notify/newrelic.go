package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const newRelicEventsURL = "https://insights-collector.newrelic.com/v1/accounts/%s/events"

// NewRelicClient sends alert events to New Relic Insights.
type NewRelicClient struct {
	accountID string
	apiKey    string
	endpoint  string
	httpClient *http.Client
}

type newRelicPayload struct {
	EventType string `json:"eventType"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
}

// NewNewRelicClient creates a NewRelicClient. accountID and apiKey must be non-empty.
func NewNewRelicClient(accountID, apiKey string) (*NewRelicClient, error) {
	if accountID == "" {
		return nil, fmt.Errorf("newrelic: account ID must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("newrelic: API key must not be empty")
	}
	return &NewRelicClient{
		accountID:  accountID,
		apiKey:     apiKey,
		endpoint:   fmt.Sprintf(newRelicEventsURL, accountID),
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert message to New Relic as a custom event.
func (c *NewRelicClient) Send(message string) error {
	payload := newRelicPayload{
		EventType: "VaultPulseAlert",
		Message:   message,
		Severity:  "critical",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("newrelic: failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}
