package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyClient sends alerts to PagerDuty via the Events API v2.
type PagerDutyClient struct {
	integrationKey string
	httpClient     *http.Client
	eventsURL      string
}

type pagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     pagerDutyInner    `json:"payload"`
}

type pagerDutyInner struct {
	Summary   string            `json:"summary"`
	Source    string            `json:"source"`
	Severity  string            `json:"severity"`
	Timestamp string            `json:"timestamp"`
	CustomDetails map[string]string `json:"custom_details,omitempty"`
}

// NewPagerDutyClient creates a new PagerDutyClient with the given integration key.
func NewPagerDutyClient(integrationKey string) (*PagerDutyClient, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty: integration key must not be empty")
	}
	return &PagerDutyClient{
		integrationKey: integrationKey,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		eventsURL:      pagerDutyEventsURL,
	}, nil
}

// Send triggers a PagerDuty alert with the given message and severity.
// severity should be one of: critical, error, warning, info.
func (c *PagerDutyClient) Send(message, severity string, details map[string]string) error {
	payload := pagerDutyPayload{
		RoutingKey:  c.integrationKey,
		EventAction: "trigger",
		Payload: pagerDutyInner{
			Summary:       message,
			Source:        "vaultpulse",
			Severity:      severity,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			CustomDetails: details,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pagerduty: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.eventsURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pagerduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
