package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewChatworkClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewChatworkClient("", "12345")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestNewChatworkClient_EmptyRoomID_ReturnsError(t *testing.T) {
	_, err := NewChatworkClient("tok", "")
	if err == nil {
		t.Fatal("expected error for empty room ID, got nil")
	}
}

func TestNewChatworkClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewChatworkClient("tok", "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestChatworkClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotBody string
	var gotToken string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		gotBody = r.FormValue("body")
		gotToken = r.Header.Get("X-ChatWorkToken")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewChatworkClient("mytoken", "99")
	c.baseURL = ts.URL

	if err := c.Send("hello chatwork"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody != "hello chatwork" {
		t.Errorf("expected body 'hello chatwork', got %q", gotBody)
	}
	if gotToken != "mytoken" {
		t.Errorf("expected token 'mytoken', got %q", gotToken)
	}
}

func TestChatworkClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewChatworkClient("tok", "1")
	c.baseURL = ts.URL

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error on non-OK status, got nil")
	}
}

func TestChatworkClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewChatworkClient("tok", "1")
	c.baseURL = "http://127.0.0.1:0"

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
