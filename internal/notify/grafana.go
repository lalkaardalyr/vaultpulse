package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GrafanaClient sends alert annotations to a Grafana instance via its HTTP API.
type GrafanaClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type grafanaAnnotation struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

// NewGrafanaClient creates a new GrafanaClient.
// baseURL must be the root URL of the Grafana instance (e.g. https://grafana.example.com).
// apiKey must be a valid Grafana service account token or API key.
func NewGrafanaClient(baseURL, apiKey string) (*GrafanaClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("grafana: base URL must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("grafana: API key must not be empty")
	}
	return &GrafanaClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}, nil
}

// Send posts an annotation to Grafana with the alert message as text.
func (g *GrafanaClient) Send(message string) error {
	payload := grafanaAnnotation{
		Text: message,
		Tags: []string{"vaultpulse", "secret-expiry"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("grafana: failed to marshal payload: %w", err)
	}

	url := g.baseURL + "/api/annotations"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("grafana: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("grafana: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
