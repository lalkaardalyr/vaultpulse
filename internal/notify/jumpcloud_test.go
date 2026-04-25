package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewJumpCloudClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewJumpCloudClient("", "org-123")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestNewJumpCloudClient_EmptyOrgID_ReturnsError(t *testing.T) {
	_, err := NewJumpCloudClient("key-abc", "")
	if err == nil {
		t.Fatal("expected error for empty org ID")
	}
}

func TestNewJumpCloudClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewJumpCloudClient("key-abc", "org-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestJumpCloudClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewJumpCloudClient("key-abc", "org-123")
	c.endpoint = ts.URL

	if err := c.Send("vault alert: secret expiring"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message"] != "vault alert: secret expiring" {
		t.Errorf("unexpected message: %v", received["message"])
	}
}

func TestJumpCloudClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewJumpCloudClient("key-abc", "org-123")
	c.endpoint = ts.URL

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestJumpCloudClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewJumpCloudClient("key-abc", "org-123")
	c.endpoint = "http://127.0.0.1:0"

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
