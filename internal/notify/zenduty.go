package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ZendutyClient sends alert notifications to Zenduty via its Events API.
type ZendutyClient struct {
	integrationKey string
	endpoint       string
	httpClient     *http.Client
}

type zendutyPayload struct {
	AlertType string `json:"alert_type"`
	Message   string `json:"message"`
	Summary   string `json:"summary"`
}

// NewZendutyClient constructs a ZendutyClient using the provided integration key.
// Returns an error if the integration key is empty.
func NewZendutyClient(integrationKey string) (*ZendutyClient, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("zenduty: integration key must not be empty")
	}
	return &ZendutyClient{
		integrationKey: integrationKey,
		endpoint:       "https://events.zenduty.com/api/events/",
		httpClient:     &http.Client{},
	}, nil
}

// Send delivers the message to Zenduty as a critical alert event.
func (z *ZendutyClient) Send(message string) error {
	payload := zendutyPayload{
		AlertType: "critical",
		Message:   message,
		Summary:   message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zenduty: failed to marshal payload: %w", err)
	}

	url := z.endpoint + z.integrationKey + "/"
	resp, err := z.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zenduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
