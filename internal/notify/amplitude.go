package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const amplitudeEndpoint = "https://api2.amplitude.com/2/httpapi"

// AmplitudeClient sends VaultPulse alert events to Amplitude Analytics.
type AmplitudeClient struct {
	apiKey   string
	endpoint string
	httpClient *http.Client
}

// NewAmplitudeClient creates a new AmplitudeClient.
// apiKey must not be empty.
func NewAmplitudeClient(apiKey string) (*AmplitudeClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("amplitude: api key must not be empty")
	}
	return &AmplitudeClient{
		apiKey:     apiKey,
		endpoint:   amplitudeEndpoint,
		httpClient: &http.Client{},
	}, nil
}

type amplitudePayload struct {
	APIKey string           `json:"api_key"`
	Events []amplitudeEvent `json:"events"`
}

type amplitudeEvent struct {
	UserID    string `json:"user_id"`
	EventType string `json:"event_type"`
	EventProperties map[string]string `json:"event_properties"`
}

// Send delivers the alert message as an Amplitude event.
func (c *AmplitudeClient) Send(message string) error {
	payload := amplitudePayload{
		APIKey: c.apiKey,
		Events: []amplitudeEvent{
			{
				UserID:    "vaultpulse",
				EventType: "vault_secret_alert",
				EventProperties: map[string]string{
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("amplitude: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
