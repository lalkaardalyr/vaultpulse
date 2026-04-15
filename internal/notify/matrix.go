package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// MatrixClient sends alerts to a Matrix room via the Client-Server API.
type MatrixClient struct {
	homeserver string
	accessToken string
	roomID      string
	httpClient  *http.Client
}

type matrixPayload struct {
	MsgType string `json:"msgtype"`
	Body    string `json:"body"`
}

// NewMatrixClient creates a new MatrixClient.
// homeserver must be the base URL (e.g. https://matrix.org),
// roomID must be the full room ID (e.g. !abc:matrix.org).
func NewMatrixClient(homeserver, accessToken, roomID string) (*MatrixClient, error) {
	if strings.TrimSpace(homeserver) == "" {
		return nil, fmt.Errorf("matrix: homeserver URL must not be empty")
	}
	if strings.TrimSpace(accessToken) == "" {
		return nil, fmt.Errorf("matrix: access token must not be empty")
	}
	if strings.TrimSpace(roomID) == "" {
		return nil, fmt.Errorf("matrix: room ID must not be empty")
	}
	return &MatrixClient{
		homeserver:  strings.TrimRight(homeserver, "/"),
		accessToken: accessToken,
		roomID:      roomID,
		httpClient:  &http.Client{},
	}, nil
}

// Send posts a plain-text message to the configured Matrix room.
func (c *MatrixClient) Send(message string) error {
	payload := matrixPayload{
		MsgType: "m.text",
		Body:    message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("matrix: failed to marshal payload: %w", err)
	}

	encoded := strings.NewReplacer("!", "%21", ":", "%3A").Replace(c.roomID)
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message", c.homeserver, encoded)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("matrix: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}
