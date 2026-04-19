package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHipChatClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewHipChatClient("", "token")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewHipChatClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewHipChatClient("https://example.com", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewHipChatClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewHipChatClient("https://example.com", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestHipChatClient_Send_PostsCorrectPayload(t *testing.T) {
	var received hipChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c, _ := NewHipChatClient(ts.URL, "mytoken")
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Message != "test alert" {
		t.Errorf("expected message 'test alert', got %q", received.Message)
	}
	if !received.Notify {
		t.Error("expected notify=true")
	}
}

func TestHipChatClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewHipChatClient(ts.URL, "tok")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestHipChatClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewHipChatClient("http://127.0.0.1:0", "tok")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
