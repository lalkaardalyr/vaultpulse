package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOpsGenieClient_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewOpsGenieClient("")
	if err == nil {
		t.Fatal("expected error for empty api key, got nil")
	}
}

func TestNewOpsGenieClient_ValidKey_ReturnsClient(t *testing.T) {
	c, err := NewOpsGenieClient("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestOpsGenieClient_Send_PostsCorrectPayload(t *testing.T) {
	var receivedAuth string
	var receivedContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	c, _ := NewOpsGenieClient("my-key")
	c.httpURL = ts.URL

	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedAuth != "GenieKey my-key" {
		t.Errorf("expected Authorization header 'GenieKey my-key', got %q", receivedAuth)
	}
	if receivedContentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", receivedContentType)
	}
}

func TestOpsGenieClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewOpsGenieClient("my-key")
	c.httpURL = ts.URL

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestOpsGenieClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewOpsGenieClient("my-key")
	c.httpURL = "http://127.0.0.1:1" // nothing listening here

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
