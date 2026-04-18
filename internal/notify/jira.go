package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// JiraClient sends alert notifications as Jira issues via the REST API.
type JiraClient struct {
	baseURL  string
	email    string
	apiToken string
	project  string
	client   *http.Client
}

type jiraIssuePayload struct {
	Fields jiraFields `json:"fields"`
}

type jiraFields struct {
	Project     jiraKey `json:"project"`
	Summary     string  `json:"summary"`
	Description string  `json:"description"`
	IssueType   jiraKey `json:"issuetype"`
}

type jiraKey struct {
	Key string `json:"key"`
}

// NewJiraClient constructs a JiraClient. baseURL should be the Jira instance
// root (e.g. https://yourorg.atlassian.net), project is the project key.
func NewJiraClient(baseURL, email, apiToken, project string) (*JiraClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("jira: baseURL must not be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("jira: email must not be empty")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("jira: apiToken must not be empty")
	}
	if project == "" {
		return nil, fmt.Errorf("jira: project must not be empty")
	}
	return &JiraClient{
		baseURL:  baseURL,
		email:    email,
		apiToken: apiToken,
		project:  project,
		client:   &http.Client{},
	}, nil
}

// Send creates a Jira issue with the alert message as the summary.
func (j *JiraClient) Send(message string) error {
	payload := jiraIssuePayload{
		Fields: jiraFields{
			Project:     jiraKey{Key: j.project},
			Summary:     message,
			Description: message,
			IssueType:   jiraKey{Key: "Task"},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("jira: marshal payload: %w", err)
	}
	url := j.baseURL + "/rest/api/2/issue"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jira: create request: %w", err)
	}
	req.SetBasicAuth(j.email, j.apiToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
