package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRocketChatClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewRocketChatClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewRocketChatClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := NewRocketChatClient("https://chat.example.com/hooks/TOKEN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestRocketChatClient_Send_PostsCorrectPayload(t *testing.T) {
	var received rocketChatPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewRocketChatClient(server.URL)
	if err := client.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text != "test alert" {
		t.Errorf("expected text %q, got %q", "test alert", received.Text)
	}
}

func TestRocketChatClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewRocketChatClient(server.URL)
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestRocketChatClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewRocketChatClient("http://127.0.0.1:0/hook")
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
