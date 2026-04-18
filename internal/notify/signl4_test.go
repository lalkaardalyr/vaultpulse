package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSIGNL4Client_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewSIGNL4Client("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewSIGNL4Client_ValidURL_ReturnsClient(t *testing.T) {
	c, err := NewSIGNL4Client("https://connect.signl4.com/webhook/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSIGNL4Client_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := NewSIGNL4Client(server.URL)
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["Message"] != "test alert" {
		t.Errorf("expected message 'test alert', got %v", received["Message"])
	}
	if received["Title"] != "VaultPulse Alert" {
		t.Errorf("expected title 'VaultPulse Alert', got %v", received["Title"])
	}
}

func TestSIGNL4Client_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c, _ := NewSIGNL4Client(server.URL)
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSIGNL4Client_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewSIGNL4Client("http://127.0.0.1:0")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
