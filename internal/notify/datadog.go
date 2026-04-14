package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const datadogEventsURL = "https://api.datadoghq.com/api/v1/events"

// DatadogClient sends alert events to the Datadog Events API.
type DatadogClient struct {
	apiKey  string
	endpoint string
	httpClient *http.Client
}

type datadogPayload struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	AlertType string `json:"alert_type"`
	Tags  []string `json:"tags,omitempty"`
}

// NewDatadogClient creates a new DatadogClient. Returns an error if apiKey is empty.
func NewDatadogClient(apiKey string) (*DatadogClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("datadog: api key must not be empty")
	}
	return &DatadogClient{
		apiKey:     apiKey,
		endpoint:   datadogEventsURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts a message to the Datadog Events API as a warning alert.
func (d *DatadogClient) Send(message string) error {
	payload := datadogPayload{
		Title:     "VaultPulse Secret Expiry Alert",
		Text:      message,
		AlertType: "warning",
		Tags:      []string{"source:vaultpulse"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("datadog: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("datadog: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
