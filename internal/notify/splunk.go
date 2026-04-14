package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SplunkClient sends alert messages to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkClient struct {
	url   string
	token string
	hc    *http.Client
}

type splunkEvent struct {
	Time       int64          `json:"time"`
	SourceType string         `json:"sourcetype"`
	Event      map[string]any `json:"event"`
}

// NewSplunkClient creates a new SplunkClient targeting the given HEC URL with
// the provided authentication token. Returns an error if either argument is empty.
func NewSplunkClient(hecURL, token string) (*SplunkClient, error) {
	if hecURL == "" {
		return nil, fmt.Errorf("splunk: HEC URL must not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("splunk: token must not be empty")
	}
	return &SplunkClient{
		url:   hecURL,
		token: token,
		hc:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts the message to the Splunk HEC endpoint as a structured event.
func (s *SplunkClient) Send(message string) error {
	payload := splunkEvent{
		Time:       time.Now().Unix(),
		SourceType: "vaultpulse",
		Event: map[string]any{
			"message": message,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("splunk: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("splunk: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Splunk "+s.token)

	resp, err := s.hc.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
