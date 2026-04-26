package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const circleCIBaseURL = "https://circleci.com/api/v2"

// CircleCIClient triggers a CircleCI pipeline as an alert action.
type CircleCIClient struct {
	token      string
	project    string
	baseURL    string
	httpClient *http.Client
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
		baseURL:    circleCIBaseURL,
		httpClient: &http.Client{},
	}, nil
}

// Send triggers a CircleCI pipeline with the alert message as a parameter.
func (c *CircleCIClient) Send(message string) error {
	url := fmt.Sprintf("%s/project/%s/pipeline", c.baseURL, c.project)
	payload := map[string]interface{}{
		"parameters": map[string]string{"alert_message": message},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("circleci: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("circleci: failed to create request: %w", err)
	}
	req.Header.Set("Circle-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("circleci: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("circleci: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
