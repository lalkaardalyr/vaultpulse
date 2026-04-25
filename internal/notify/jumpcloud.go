package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const jumpCloudDefaultEndpoint = "https://api.jumpcloud.com/insights/directory/v1/events"

// JumpCloudClient sends alert events to the JumpCloud Insights API.
type JumpCloudClient struct {
	apiKey     string
	orgID      string
	endpoint   string
	httpClient *http.Client
}

type jumpCloudPayload struct {
	Message string `json:"message"`
	OrgID   string `json:"orgId"`
}

// NewJumpCloudClient creates a new JumpCloudClient.
// Returns an error if apiKey or orgID is empty.
func NewJumpCloudClient(apiKey, orgID string) (*JumpCloudClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("jumpcloud: API key must not be empty")
	}
	if orgID == "" {
		return nil, fmt.Errorf("jumpcloud: org ID must not be empty")
	}
	return &JumpCloudClient{
		apiKey:     apiKey,
		orgID:      orgID,
		endpoint:   jumpCloudDefaultEndpoint,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert event to the JumpCloud Insights API.
func (c *JumpCloudClient) Send(message string) error {
	payload := jumpCloudPayload{
		Message: message,
		OrgID:   c.orgID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("jumpcloud: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jumpcloud: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("jumpcloud: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jumpcloud: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
