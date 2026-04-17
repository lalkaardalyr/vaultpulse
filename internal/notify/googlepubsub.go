package notify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const pubsubEndpoint = "https://pubsub.googleapis.com/v1/projects/%s/topics/%s:publish"

// GooglePubSubClient sends alert messages to a Google Cloud Pub/Sub topic.
type GooglePubSubClient struct {
	projectID string
	topicID   string
	apiKey    string
	httpClient *http.Client
	endpoint  string
}

type pubsubMessage struct {
	Messages []pubsubEntry `json:"messages"`
}

type pubsubEntry struct {
	Data string `json:"data"`
}

// NewGooglePubSubClient creates a new GooglePubSubClient.
func NewGooglePubSubClient(projectID, topicID, apiKey string) (*GooglePubSubClient, error) {
	if projectID == "" {
		return nil, errors.New("googlepubsub: project ID must not be empty")
	}
	if topicID == "" {
		return nil, errors.New("googlepubsub: topic ID must not be empty")
	}
	if apiKey == "" {
		return nil, errors.New("googlepubsub: API key must not be empty")
	}
	return &GooglePubSubClient{
		projectID:  projectID,
		topicID:    topicID,
		apiKey:     apiKey,
		httpClient: &http.Client{},
		endpoint:   fmt.Sprintf(pubsubEndpoint, projectID, topicID),
	}, nil
}

// Send publishes the message to the configured Pub/Sub topic.
func (c *GooglePubSubClient) Send(_ context.Context, message string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(message))
	payload := pubsubMessage{Messages: []pubsubEntry{{Data: encoded}}}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlepubsub: marshal error: %w", err)
	}
	url := c.endpoint + "?key=" + c.apiKey
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlepubsub: request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("googlepubsub: unexpected status %d", resp.StatusCode)
	}
	return nil
}
