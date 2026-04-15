package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SignalSciencesClient sends alert notifications to Signal Sciences (Fastly Next-Gen WAF)
// via its custom event ingestion API.
type SignalSciencesClient struct {
	endpoint string
	apiUser  string
	apiToken string
	corpName string
	client   *http.Client
}

type signalSciencesPayload struct {
	Event   string            `json:"event"`
	Message string            `json:"message"`
	Tags    []string          `json:"tags"`
	Meta    map[string]string `json:"meta,omitempty"`
}

// NewSignalSciencesClient creates a new SignalSciencesClient.
// apiUser and apiToken are required for authentication; corpName identifies the corp.
func NewSignalSciencesClient(apiUser, apiToken, corpName string) (*SignalSciencesClient, error) {
	if apiUser == "" {
		return nil, fmt.Errorf("signalsciences: api user must not be empty")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("signalsciences: api token must not be empty")
	}
	if corpName == "" {
		return nil, fmt.Errorf("signalsciences: corp name must not be empty")
	}
	return &SignalSciencesClient{
		endpoint: fmt.Sprintf("https://dashboard.signalsciences.net/api/v0/corps/%s/events", corpName),
		apiUser:  apiUser,
		apiToken: apiToken,
		corpName: corpName,
		client:   &http.Client{},
	}, nil
}

// Send delivers the alert message to Signal Sciences as a custom event.
func (c *SignalSciencesClient) Send(message string) error {
	payload := signalSciencesPayload{
		Event:   "vaultpulse.alert",
		Message: message,
		Tags:    []string{"vaultpulse", "secret-expiry"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalsciences: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalsciences: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-user", c.apiUser)
	req.Header.Set("x-api-token", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalsciences: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalsciences: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
