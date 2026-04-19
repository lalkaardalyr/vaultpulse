package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRevoltClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewRevoltClient("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewRevoltClient_ValidURL_ReturnsClient(t *testing.T) {
	c, err := NewRevoltClient("https://revolt.example.com/webhook/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestRevoltClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := NewRevoltClient(server.URL)
	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["content"] != "vault secret expiring soon" {
		t.Errorf("unexpected content: %q", received["content"])
	}
}

func TestRevoltClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c, _ := NewRevoltClient(server.URL)
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestRevoltClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewRevoltClient("http://127.0.0.1:0")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
