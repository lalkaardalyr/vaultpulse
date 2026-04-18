package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAmplitudeClient_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewAmplitudeClient("")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewAmplitudeClient_ValidKey_ReturnsClient(t *testing.T) {
	c, err := NewAmplitudeClient("test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestAmplitudeClient_Send_PostsCorrectPayload(t *testing.T) {
	var received amplitudePayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewAmplitudeClient("my-key")
	c.endpoint = ts.URL

	if err := c.Send("secret expires soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.APIKey != "my-key" {
		t.Errorf("expected api key %q, got %q", "my-key", received.APIKey)
	}
	if len(received.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received.Events))
	}
	if received.Events[0].EventType != "vault_secret_alert" {
		t.Errorf("unexpected event type: %s", received.Events[0].EventType)
	}
	if received.Events[0].EventProperties["message"] != "secret expires soon" {
		t.Errorf("unexpected message: %s", received.Events[0].EventProperties["message"])
	}
}

func TestAmplitudeClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewAmplitudeClient("key")
	c.endpoint = ts.URL

	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestAmplitudeClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewAmplitudeClient("key")
	c.endpoint = "http://127.0.0.1:0"

	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
