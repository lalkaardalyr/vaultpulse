package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNewRelicClient_EmptyAccountID_ReturnsError(t *testing.T) {
	_, err := NewNewRelicClient("", "key")
	if err == nil {
		t.Fatal("expected error for empty account ID")
	}
}

func TestNewNewRelicClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewNewRelicClient("12345", "")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestNewNewRelicClient_ValidConfig_ReturnsClient(t *testing.T) {
	client, err := NewNewRelicClient("12345", "secret-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewRelicClient_Send_PostsCorrectPayload(t *testing.T) {
	var received newRelicPayload
	var gotKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-Insert-Key")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewNewRelicClient("12345", "test-key")
	client.endpoint = server.URL
	client.httpClient = server.Client()

	if err := client.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Message != "test alert" {
		t.Errorf("expected message 'test alert', got %q", received.Message)
	}
	if received.EventType != "VaultPulseAlert" {
		t.Errorf("expected eventType 'VaultPulseAlert', got %q", received.EventType)
	}
	if gotKey != "test-key" {
		t.Errorf("expected X-Insert-Key 'test-key', got %q", gotKey)
	}
}

func TestNewRelicClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client, _ := NewNewRelicClient("12345", "test-key")
	client.endpoint = server.URL
	client.httpClient = server.Client()

	if err := client.Send("alert"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestNewRelicClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewNewRelicClient("12345", "test-key")
	client.endpoint = "http://127.0.0.1:0/events"

	if err := client.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
