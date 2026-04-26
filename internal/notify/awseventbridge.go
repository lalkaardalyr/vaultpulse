// Package notify provides notification sender implementations for various
// platforms and services. Each client implements the Sender interface.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EventBridgeClient sends alert events to AWS EventBridge via the HTTP API.
type EventBridgeClient struct {
	endpointURL string
	source      string
	detailType  string
	httpClient  *http.Client
}

// NewEventBridgeClient creates a new EventBridgeClient.
// endpointURL is the EventBridge endpoint, source and detailType label the event.
func NewEventBridgeClient(endpointURL, source, detailType string) (*EventBridgeClient, error) {
	if endpointURL == "" {
		return nil, fmt.Errorf("eventbridge: endpoint URL must not be empty")
	}
	if source == "" {
		return nil, fmt.Errorf("eventbridge: source must not be empty")
	}
	if detailType == "" {
		return nil, fmt.Errorf("eventbridge: detail type must not be empty")
	}
	return &EventBridgeClient{
		endpointURL: endpointURL,
		source:      source,
		detailType:  detailType,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts an EventBridge-compatible event payload with the alert message.
func (c *EventBridgeClient) Send(message string) error {
	payload := map[string]interface{}{
		"Source":     c.source,
		"DetailType": c.detailType,
		"Detail":     map[string]string{"message": message},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("eventbridge: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.endpointURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("eventbridge: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("eventbridge: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
