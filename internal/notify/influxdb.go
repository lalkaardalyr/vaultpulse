package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// InfluxDBClient sends alert events to InfluxDB as line protocol data
// via the /api/v2/write endpoint.
type InfluxDBClient struct {
	url    string
	token  string
	org    string
	bucket string
	httpClient *http.Client
}

// NewInfluxDBClient creates a new InfluxDBClient.
func NewInfluxDBClient(url, token, org, bucket string) (*InfluxDBClient, error) {
	if url == "" {
		return nil, errors.New("influxdb: url must not be empty")
	}
	if token == "" {
		return nil, errors.New("influxdb: token must not be empty")
	}
	if org == "" {
		return nil, errors.New("influxdb: org must not be empty")
	}
	if bucket == "" {
		return nil, errors.New("influxdb: bucket must not be empty")
	}
	return &InfluxDBClient{
		url:        url,
		token:      token,
		org:        org,
		bucket:     bucket,
		httpClient: &http.Client{},
	}, nil
}

// Send writes an alert event to InfluxDB.
func (c *InfluxDBClient) Send(ctx context.Context, message string) error {
	payload := map[string]interface{}{
		"measurement": "vaultpulse_alert",
		"fields":      map[string]interface{}{"message": message},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("influxdb: failed to marshal payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v2/write?org=%s&bucket=%s&precision=s", c.url, c.org, c.bucket)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("influxdb: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("influxdb: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("influxdb: unexpected status %d", resp.StatusCode)
	}
	return nil
}
