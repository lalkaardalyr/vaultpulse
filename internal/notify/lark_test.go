package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewLarkClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := notify.NewLarkClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewLarkClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := notify.NewLarkClient("https://open.larksuite.com/open-apis/bot/v2/hook/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestLarkClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client, err := notify.NewLarkClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("test alert message"); err != nil {
		t.Fatalf("Send returned unexpected error: %v", err)
	}

	if received["msg_type"] != "text" {
		t.Errorf("expected msg_type=text, got %v", received["msg_type"])
	}

	content, ok := received["content"].(map[string]interface{})
	if !ok {
		t.Fatal("expected content to be a map")
	}
	if content["text"] != "test alert message" {
		t.Errorf("expected text='test alert message', got %v", content["text"])
	}
}

func TestLarkClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client, err := notify.NewLarkClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("alert"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestLarkClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, err := notify.NewLarkClient("http://127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
