package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// StatuspageClient sends incident notifications to Atlassian Statuspage.
type StatuspageClient struct {
	apiKey   string
	pageID   string
	endpoint string
	httpClient *http.Client
}

type statuspageIncident struct {
	Incident statuspageBody `json:"incident"`
}

type statuspageBody struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Body   string `json:"body"`
}

// NewStatuspageClient creates a new Statuspage notification client.
func NewStatuspageClient(apiKey, pageID string) (*StatuspageClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("statuspage: api key must not be empty")
	}
	if pageID == "" {
		return nil, fmt.Errorf("statuspage: page ID must not be empty")
	}
	return &StatuspageClient{
		apiKey:     apiKey,
		pageID:     pageID,
		endpoint:   "https://api.statuspage.io/v1",
		httpClient: &http.Client{},
	}, nil
}

// Send posts an incident to Statuspage.
func (c *StatuspageClient) Send(msg string) error {
	payload := statuspageIncident{
		Incident: statuspageBody{
			Name:   "VaultPulse Alert",
			Status: "investigating",
			Body:   msg,
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("statuspage: marshal error: %w", err)
	}
	url := fmt.Sprintf("%s/pages/%s/incidents", c.endpoint, c.pageID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("statuspage: request error: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: send error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status %d", resp.StatusCode)
	}
	return nil
}
