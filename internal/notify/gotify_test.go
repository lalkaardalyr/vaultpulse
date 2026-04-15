package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGotifyClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewGotifyClient("", "token", 5)
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewGotifyClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewGotifyClient("https://gotify.example.com", "", 5)
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestNewGotifyClient_ValidConfig_ReturnsClient(t *testing.T) {
	client, err := NewGotifyClient("https://gotify.example.com", "token", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGotifyClient_Send_PostsCorrectPayload(t *testing.T) {
	var received gotifyPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewGotifyClient(server.URL, "testtoken", 7)
	if err := client.Send("test alert message"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Message != "test alert message" {
		t.Errorf("expected message %q, got %q", "test alert message", received.Message)
	}
	if received.Priority != 7 {
		t.Errorf("expected priority 7, got %d", received.Priority)
	}
	if received.Title != "VaultPulse Alert" {
		t.Errorf("expected title %q, got %q", "VaultPulse Alert", received.Title)
	}
}

func TestGotifyClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client, _ := NewGotifyClient(server.URL, "badtoken", 5)
	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestGotifyClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewGotifyClient("http://127.0.0.1:0", "token", 5)
	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
