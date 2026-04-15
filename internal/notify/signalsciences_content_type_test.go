package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignalSciencesClient_Send_PostsJSONContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, err := NewSignalSciencesClient("user@example.com", "mytoken", "mycorp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.endpoint = ts.URL

	if err := c.Send("content-type check"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
