package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPushoverClient_Send_PostsJSONContentType(t *testing.T) {
	var gotContentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewPushoverClient("tok", "usr")
	c.apiURL = ts.URL

	if err := c.Send("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", gotContentType)
	}
}
