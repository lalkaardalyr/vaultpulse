package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewEventBridgeClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewEventBridgeClient("", "src", "detail")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewEventBridgeClient_EmptySource_ReturnsError(t *testing.T) {
	_, err := NewEventBridgeClient("http://example.com", "", "detail")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNewEventBridgeClient_EmptyDetailType_ReturnsError(t *testing.T) {
	_, err := NewEventBridgeClient("http://example.com", "src", "")
	if err == nil {
		t.Fatal("expected error for empty detail type")
	}
}

func TestNewEventBridgeClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewEventBridgeClient("http://example.com", "vaultpulse", "SecretExpiry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestEventBridgeClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewEventBridgeClient(ts.URL, "vaultpulse", "SecretExpiry")
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["Source"] != "vaultpulse" {
		t.Errorf("expected source vaultpulse, got %v", received["Source"])
	}
	if received["DetailType"] != "SecretExpiry" {
		t.Errorf("expected DetailType SecretExpiry, got %v", received["DetailType"])
	}
}

func TestEventBridgeClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewEventBridgeClient(ts.URL, "vaultpulse", "SecretExpiry")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestEventBridgeClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewEventBridgeClient("http://127.0.0.1:0", "vaultpulse", "SecretExpiry")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
