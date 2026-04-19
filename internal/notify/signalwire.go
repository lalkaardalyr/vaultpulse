package notify

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SignalWireClient sends SMS alerts via the SignalWire REST API.
type SignalWireClient struct {
	projectID string
	authToken string
	spaceURL  string
	from      string
	to        string
	httpClient *http.Client
}

// NewSignalWireClient constructs a SignalWireClient.
func NewSignalWireClient(projectID, authToken, spaceURL, from, to string) (*SignalWireClient, error) {
	if projectID == "" {
		return nil, errors.New("signalwire: project ID must not be empty")
	}
	if authToken == "" {
		return nil, errors.New("signalwire: auth token must not be empty")
	}
	if spaceURL == "" {
		return nil, errors.New("signalwire: space URL must not be empty")
	}
	if from == "" {
		return nil, errors.New("signalwire: from number must not be empty")
	}
	if to == "" {
		return nil, errors.New("signalwire: to number must not be empty")
	}
	return &SignalWireClient{
		projectID:  projectID,
		authToken:  authToken,
		spaceURL:   strings.TrimRight(spaceURL, "/"),
		from:       from,
		to:         to,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the message via SignalWire's Messages API.
func (c *SignalWireClient) Send(message string) error {
	endpoint := fmt.Sprintf("%s/api/laml/2010-04-01/Accounts/%s/Messages.json",
		c.spaceURL, c.projectID)

	form := url.Values{}
	form.Set("From", c.from)
	form.Set("To", c.to)
	form.Set("Body", message)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("signalwire: failed to build request: %w", err)
	}
	req.SetBasicAuth(c.projectID, c.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("signalwire: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalwire: unexpected status %d", resp.StatusCode)
	}
	return nil
}
