package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMatrixClient_EmptyHomeserver_ReturnsError(t *testing.T) {
	_, err := NewMatrixClient("", "token", "!room:matrix.org")
	if err == nil {
		t.Fatal("expected error for empty homeserver")
	}
}

func TestNewMatrixClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewMatrixClient("https://matrix.org", "", "!room:matrix.org")
	if err == nil {
		t.Fatal("expected error for empty access token")
	}
}

func TestNewMatrixClient_EmptyRoomID_ReturnsError(t *testing.T) {
	_, err := NewMatrixClient("https://matrix.org", "token", "")
	if err == nil {
		t.Fatal("expected error for empty room ID")
	}
}

func TestNewMatrixClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewMatrixClient("https://matrix.org", "token", "!room:matrix.org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestMatrixClient_Send_PostsCorrectPayload(t *testing.T) {
	var captured []byte
	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = io.ReadAll(r.Body)
		capturedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer server.Close()

	c, _ := NewMatrixClient(server.URL, "mytoken", "!room:matrix.org")
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload matrixPayload
	if err := json.Unmarshal(captured, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload.Body != "test alert" {
		t.Errorf("expected body 'test alert', got %q", payload.Body)
	}
	if payload.MsgType != "m.text" {
		t.Errorf("expected msgtype 'm.text', got %q", payload.MsgType)
	}
	if !strings.HasPrefix(capturedAuth, "Bearer ") {
		t.Errorf("expected Bearer auth header, got %q", capturedAuth)
	}
}

func TestMatrixClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	c, _ := NewMatrixClient(server.URL, "token", "!room:matrix.org")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMatrixClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewMatrixClient("http://127.0.0.1:19999", "token", "!room:matrix.org")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
