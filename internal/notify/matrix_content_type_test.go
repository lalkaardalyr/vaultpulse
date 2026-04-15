package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatrixClient_Send_PostsJSONContentType(t *testing.T) {
	var contentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$xyz"}`))
	}))
	defer server.Close()

	c, err := NewMatrixClient(server.URL, "token", "!room:matrix.org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.Send("hello"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}
