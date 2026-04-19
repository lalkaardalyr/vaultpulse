package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSignalWireClient_EmptyProjectID_ReturnsError(t *testing.T) {
	_, err := NewSignalWireClient("", "token", "https://example.signalwire.com", "+1000", "+2000")
	if err == nil {
		t.Fatal("expected error for empty project ID")
	}
}

func TestNewSignalWireClient_EmptyAuthToken_ReturnsError(t *testing.T) {
	_, err := NewSignalWireClient("proj", "", "https://example.signalwire.com", "+1000", "+2000")
	if err == nil {
		t.Fatal("expected error for empty auth token")
	}
}

func TestNewSignalWireClient_EmptySpaceURL_ReturnsError(t *testing.T) {
	_, err := NewSignalWireClient("proj", "token", "", "+1000", "+2000")
	if err == nil {
		t.Fatal("expected error for empty space URL")
	}
}

func TestNewSignalWireClient_EmptyFrom_ReturnsError(t *testing.T) {
	_, err := NewSignalWireClient("proj", "token", "https://example.signalwire.com", "", "+2000")
	if err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestNewSignalWireClient_EmptyTo_ReturnsError(t *testing.T) {
	_, err := NewSignalWireClient("proj", "token", "https://example.signalwire.com", "+1000", "")
	if err == nil {
		t.Fatal("expected error for empty to")
	}
}

func TestNewSignalWireClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewSignalWireClient("proj", "token", "https://example.signalwire.com", "+1000", "+2000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSignalWireClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		gotBody = r.FormValue("Body")
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	c, _ := NewSignalWireClient("proj", "token", server.URL, "+1000", "+2000")
	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody != "vault secret expiring soon" {
		t.Errorf("expected body %q, got %q", "vault secret expiring soon", gotBody)
	}
}

func TestSignalWireClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	c, _ := NewSignalWireClient("proj", "token", server.URL, "+1000", "+2000")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
