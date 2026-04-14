package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const opsGenieAlertsURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieClient sends alerts to OpsGenie via the REST API.
type OpsGenieClient struct {
	apiKey  string
	httpURL string // overridable for testing
	client  *http.Client
}

type opsGeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description,omitempty"`
	Priority    string            `json:"priority"`
	Details     map[string]string `json:"details,omitempty"`
}

// NewOpsGenieClient creates an OpsGenieClient using the provided API key.
// Returns an error if the key is empty.
func NewOpsGenieClient(apiKey string) (*OpsGenieClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("opsgenie: api key must not be empty")
	}
	return &OpsGenieClient{
		apiKey:  apiKey,
		httpURL: opsGenieAlertsURL,
		client:  &http.Client{},
	}, nil
}

// Send delivers the message to OpsGenie as a P1 alert.
func (o *OpsGenieClient) Send(message string) error {
	payload := opsGeniePayload{
		Message:  message,
		Priority: "P1",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.httpURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("opsgenie: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
