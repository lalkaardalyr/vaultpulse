package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const linearAPIURL = "https://api.linear.app/graphql"

// LinearClient creates Linear issues for vault alerts.
type LinearClient struct {
	apiKey  string
	teamID  string
	endpoint string
	httpClient *http.Client
}

// NewLinearClient creates a new LinearClient.
func NewLinearClient(apiKey, teamID string) (*LinearClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("linear: API key must not be empty")
	}
	if teamID == "" {
		return nil, fmt.Errorf("linear: team ID must not be empty")
	}
	return &LinearClient{
		apiKey:     apiKey,
		teamID:     teamID,
		endpoint:   linearAPIURL,
		httpClient: &http.Client{},
	}, nil
}

// Send creates a Linear issue with the given message as the title.
func (c *LinearClient) Send(msg string) error {
	query := fmt.Sprintf(`{"query":"mutation { issueCreate(input: { title: %q, teamId: %q }) { success } }"}`, msg, c.teamID)
	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBufferString(query))
	if err != nil {
		return fmt.Errorf("linear: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("linear: request error: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Data struct {
			IssueCreate struct {
				Success bool `json:"success"`
			} `json:"issueCreate"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("linear: decode response: %w", err)
	}
	if !result.Data.IssueCreate.Success {
		return fmt.Errorf("linear: issue creation failed")
	}
	return nil
}
