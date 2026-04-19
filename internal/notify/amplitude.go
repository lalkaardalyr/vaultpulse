package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const amplitudeEndpoint = "https://api2.amplitude.com/2/httpapi"

// AmplitudeClient sends events to Amplitude.
type AmplitudeClient struct {
	apiKey   string
	endpoint string
	httpClient *http.Client
}

// NewAmplitudeClient creates a new AmplitudeClient.
func NewAmplitudeClient(apiKey string) (*AmplitudeClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("amplitude: API key must not be empty")
	}
	return &AmplitudeClient{
		apiKey:     apiKey,
		endpoint:   amplitudeEndpoint,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an event to Amplitude.
func (c *AmplitudeClient) Send(message string) error {
	payload := map[string]interface{}{
		"api_key": c.apiKey,
		"events": []map[string]interface{}{
			{
				"event_type":    "vaultpulse_alert",
				"user_id":       "vaultpulse",
				"event_properties": map[string]string{
					"message": message,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("amplitude: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("amplitude: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("amplitude: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
