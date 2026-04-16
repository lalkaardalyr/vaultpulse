package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNtfyClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewNtfyClient("", "alerts")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewNtfyClient_EmptyTopic_ReturnsError(t *testing.T) {
	_, err := NewNtfyClient("https://ntfy.sh", "")
	if err == nil {
		t.Fatal("expected error for empty topic")
	}
}

func TestNewNtfyClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewNtfyClient("https://ntfy.sh", "alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNtfyClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotTopic, gotBody, gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotTopic = r.URL.Path
		buf := make([]byte, 512)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewNtfyClient(ts.URL, "vaultpulse")
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotTopic != "/vaultpulse" {
		t.Errorf("expected topic path /vaultpulse, got %s", gotTopic)
	}
	if gotBody != "test alert" {
		t.Errorf("unexpected body: %s", gotBody)
	}
}

func TestNtfyClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewNtfyClient(ts.URL, "vaultpulse")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestNtfyClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewNtfyClient("http://127.0.0.1:19999", "vaultpulse")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
