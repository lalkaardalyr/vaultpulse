package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewChatworkClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewChatworkClient("", "12345")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewChatworkClient_EmptyRoomID_ReturnsError(t *testing.T) {
	_, err := NewChatworkClient("token", "")
	if err == nil {
		t.Fatal("expected error for empty roomID")
	}
}

func TestNewChatworkClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewChatworkClient("token", "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestChatworkClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotToken, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("X-ChatWorkToken")
		r.ParseForm()
		gotBody = r.FormValue("body")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewChatworkClient("mytoken", "99")
	c.baseURL = ts.URL

	if err := c.Send("hello vault"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotToken != "mytoken" {
		t.Errorf("expected token mytoken, got %s", gotToken)
	}
	if gotBody != "hello vault" {
		t.Errorf("expected body 'hello vault', got %s", gotBody)
	}
}

func TestChatworkClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewChatworkClient("token", "99")
	c.baseURL = ts.URL

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestChatworkClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewChatworkClient("token", "99")
	c.baseURL = "http://127.0.0.1:0"

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error on unreachable server")
	}
}
