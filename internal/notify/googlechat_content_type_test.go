package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoogleChatClient_Send_PostsJSONContentType(t *testing.T) {
	var capturedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewGoogleChatClient(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("content-type check"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if capturedContentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", capturedContentType)
	}
}
