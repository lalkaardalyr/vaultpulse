package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTelegramClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewTelegramClient("", "123456")
	if err == nil {
		t.Fatal("expected error for empty bot token, got nil")
	}
}

func TestNewTelegramClient_EmptyChatID_ReturnsError(t *testing.T) {
	_, err := NewTelegramClient("bot-token", "")
	if err == nil {
		t.Fatal("expected error for empty chat ID, got nil")
	}
}

func TestNewTelegramClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewTelegramClient("bot-token", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestTelegramClient_Send_PostsCorrectPayload(t *testing.T) {
	var received telegramPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := NewTelegramClient("test-token", "chat-99")
	c.endpoint = server.URL

	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.ChatID != "chat-99" {
		t.Errorf("expected chat_id %q, got %q", "chat-99", received.ChatID)
	}
	if received.Text != "vault secret expiring soon" {
		t.Errorf("unexpected text: %q", received.Text)
	}
	if received.ParseMode != "Markdown" {
		t.Errorf("expected parse_mode Markdown, got %q", received.ParseMode)
	}
}

func TestTelegramClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	c, _ := NewTelegramClient("bad-token", "chat-99")
	c.endpoint = server.URL

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestTelegramClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewTelegramClient("token", "chat-1")
	c.endpoint = "http://127.0.0.1:0/sendMessage"

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
