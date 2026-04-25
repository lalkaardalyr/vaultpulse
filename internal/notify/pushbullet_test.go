package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPushbulletClient_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewPushbulletClient("")
	if err == nil {
		t.Fatal("expected error for empty api key, got nil")
	}
}

func TestNewPushbulletClient_ValidKey_ReturnsClient(t *testing.T) {
	c, err := NewPushbulletClient("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestPushbulletClient_Send_PostsCorrectPayload(t *testing.T) {
	var received pushbulletPayload
	var gotToken string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("Access-Token")
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewPushbulletClient("my-key")
	c.url = ts.URL

	if err := c.Send("secret expires soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Body != "secret expires soon" {
		t.Errorf("expected body %q, got %q", "secret expires soon", received.Body)
	}
	if received.Type != "note" {
		t.Errorf("expected type 'note', got %q", received.Type)
	}
	if gotToken != "my-key" {
		t.Errorf("expected Access-Token 'my-key', got %q", gotToken)
	}
}

func TestPushbulletClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewPushbulletClient("bad-key")
	c.url = ts.URL

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestPushbulletClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewPushbulletClient("key")
	c.url = "http://127.0.0.1:0"

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
