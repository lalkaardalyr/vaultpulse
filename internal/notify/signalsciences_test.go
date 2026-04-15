package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSignalSciencesClient_EmptyUser_ReturnsError(t *testing.T) {
	_, err := NewSignalSciencesClient("", "token", "mycorp")
	if err == nil {
		t.Fatal("expected error for empty api user")
	}
}

func TestNewSignalSciencesClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewSignalSciencesClient("user@example.com", "", "mycorp")
	if err == nil {
		t.Fatal("expected error for empty api token")
	}
}

func TestNewSignalSciencesClient_EmptyCorp_ReturnsError(t *testing.T) {
	_, err := NewSignalSciencesClient("user@example.com", "token", "")
	if err == nil {
		t.Fatal("expected error for empty corp name")
	}
}

func TestNewSignalSciencesClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewSignalSciencesClient("user@example.com", "token", "mycorp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSignalSciencesClient_Send_PostsCorrectPayload(t *testing.T) {
	var received signalSciencesPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		if r.Header.Get("x-api-user") == "" {
			t.Error("expected x-api-user header")
		}
		if r.Header.Get("x-api-token") == "" {
			t.Error("expected x-api-token header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewSignalSciencesClient("user@example.com", "mytoken", "mycorp")
	c.endpoint = ts.URL

	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Message != "test alert" {
		t.Errorf("expected message 'test alert', got %q", received.Message)
	}
	if received.Event != "vaultpulse.alert" {
		t.Errorf("expected event 'vaultpulse.alert', got %q", received.Event)
	}
}

func TestSignalSciencesClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewSignalSciencesClient("user@example.com", "mytoken", "mycorp")
	c.endpoint = ts.URL

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSignalSciencesClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewSignalSciencesClient("user@example.com", "mytoken", "mycorp")
	c.endpoint = "http://127.0.0.1:0/invalid"

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
