package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const circleciBaseURL = "https://circleci.com/api/v2"

// CircleCIClient triggers a CircleCI pipeline as an alert notification.
type CircleCIClient struct {
	token      string
	project    string
	baseURL    string
	httpClient *http.Client
}

type circleCIPayload struct {
	Branch     string            `json:"branch"`
	Parameters map[string]string `json:"parameters"`
}

// NewCircleCIClient creates a new CircleCIClient.
// Returns an error if token or project is empty.
func NewCircleCIClient(token, project string) (*CircleCIClient, error) {
	if token == "" {
		return nil, fmt.Errorf("circleci: token must not be empty")
	}
	if project == "" {
		return nil, fmt.Errorf("circleci: project must not be empty")
	}
	return &CircleCIClient{
		token:      token,
		project:    project,
		baseURL:    circleciBaseURL,
		httpClient: &http.Client{},
	}, nil
}

// Send triggers a CircleCI pipeline with the alert message as a parameter.
func (c *CircleCIClient) Send(message string) error {
	payload := circleCIPayload{
		Branch: "main",
		Parameters: map[string]string{
			"alert_message": message,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("circleci: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/project/%s/pipeline", c.baseURL, c.project)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("circleci: create request: %w", err)
	}
	req.Header.Set("Circle-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("circleci: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("circleci: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
