package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewBearyChatClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewBearyChatClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewBearyChatClient_ValidURL_ReturnsClient(t *testing.T) {
	c, err := NewBearyChatClient("https://hook.bearychat.com/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestBearyChatClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewBearyChatClient(ts.URL)
	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["text"] != "vault secret expiring soon" {
		t.Errorf("expected text field to match message, got %q", received["text"])
	}
	if received["notification"] != "vault secret expiring soon" {
		t.Errorf("expected notification field to match message, got %q", received["notification"])
	}
}

func TestBearyChatClient_Send_PostsCorrectContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewBearyChatClient(ts.URL)
	if err := c.Send("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBearyChatClient_Send_UsesPostMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %q", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewBearyChatClient(ts.URL)
	if err := c.Send("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBearyChatClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewBearyChatClient(ts.URL)
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestBearyChatClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewBearyChatClient("http://127.0.0.1:0")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
