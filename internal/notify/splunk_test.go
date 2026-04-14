package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewSplunkClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewSplunkClient("", "token")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewSplunkClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewSplunkClient("http://splunk.example.com/services/collector", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestNewSplunkClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewSplunkClient("http://splunk.example.com/services/collector", "mytoken")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSplunkClient_Send_PostsCorrectPayload(t *testing.T) {
	var received []byte
	var authHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewSplunkClient(ts.URL, "test-token")
	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(authHeader, "Splunk ") {
		t.Errorf("expected Authorization header to start with 'Splunk ', got %q", authHeader)
	}

	var payload splunkEvent
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}
	if payload.SourceType != "vaultpulse" {
		t.Errorf("expected sourcetype 'vaultpulse', got %q", payload.SourceType)
	}
	msg, _ := payload.Event["message"].(string)
	if msg != "vault secret expiring soon" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestSplunkClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewSplunkClient(ts.URL, "bad-token")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestSplunkClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewSplunkClient("http://127.0.0.1:19999/services/collector", "token")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
