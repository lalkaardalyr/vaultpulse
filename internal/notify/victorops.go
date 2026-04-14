package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// VictorOpsClient sends alert notifications to a VictorOps (Splunk On-Call) REST endpoint.
type VictorOpsClient struct {
	webhookURL string
	httpClient *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	MonitoringTool    string `json:"monitoring_tool"`
}

// NewVictorOpsClient constructs a VictorOpsClient for the given REST endpoint URL.
// Returns an error if the URL is empty.
func NewVictorOpsClient(webhookURL string) (*VictorOpsClient, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("victorops: webhook URL must not be empty")
	}
	return &VictorOpsClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the message to VictorOps as a CRITICAL alert.
func (c *VictorOpsClient) Send(message string) error {
	payload := victorOpsPayload{
		MessageType:       "CRITICAL",
		EntityDisplayName: "VaultPulse Secret Expiry Alert",
		StateMessage:      message,
		MonitoringTool:    "vaultpulse",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
