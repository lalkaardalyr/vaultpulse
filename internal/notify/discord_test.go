package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDiscordClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewDiscordClient("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}

func TestNewDiscordClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := NewDiscordClient("https://discord.com/api/webhooks/123/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestDiscordClient_Send_PostsCorrectPayload(t *testing.T) {
	var received discordPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewDiscordClient(server.URL)
	if err := client.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Content != "test alert" {
		t.Errorf("expected content %q, got %q", "test alert", received.Content)
	}
}

// TestDiscordClient_Send_PostsJSONContentType verifies that the client sends
// requests with the correct Content-Type header for Discord's API.
func TestDiscordClient_Send_PostsJSONContentType(t *testing.T) {
	var contentType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewDiscordClient(server.URL)
	if err := client.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contentType != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", contentType)
	}
}

func TestDiscordClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client, _ := NewDiscordClient(server.URL)
	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestDiscordClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewDiscordClient("http://127.0.0.1:0/webhook")
	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
