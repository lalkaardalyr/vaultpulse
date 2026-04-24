package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const statuspageBaseURL = "https://api.statuspage.io/v1"

// StatuspageClient sends incident updates to Atlassian Statuspage.
type StatuspageClient struct {
	apiKey  string
	pageID  string
	baseURL string
	httpClient *http.Client
}

// NewStatuspageClient creates a new StatuspageClient.
func NewStatuspageClient(apiKey, pageID string) (*StatuspageClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("statuspage: API key must not be empty")
	}
	if pageID == "" {
		return nil, fmt.Errorf("statuspage: page ID must not be empty")
	}
	return &StatuspageClient{
		apiKey:     apiKey,
		pageID:     pageID,
		baseURL:    statuspageBaseURL,
		httpClient: &http.Client{},
	}, nil
}

type statuspageIncident struct {
	Incident statuspageIncidentBody `json:"incident"`
}

type statuspageIncidentBody struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

// Send creates a new incident on Statuspage with the given message.
func (c *StatuspageClient) Send(message string) error {
	payload := statuspageIncident{
		Incident: statuspageIncidentBody{
			Name: "VaultPulse Alert",
			Body: message,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("statuspage: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/pages/%s/incidents", c.baseURL, c.pageID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("statuspage: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
