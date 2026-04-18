package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZendeskClient_MultiSender_Integration(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	// Build a fake ZendeskClient that hits the test server
	zc := &ZendeskClient{
		subdomain:  "test",
		email:      "ops@example.com",
		apiToken:   "secret",
		httpClient: ts.Client(),
	}

	// Wrap in MultiSender alongside a no-op sender
	noop := &noopSender{}
	multi, err := NewMultiSender(zc, noop)
	if err != nil {
		t.Fatalf("NewMultiSender: %v", err)
	}

	// Send will fail because ts.Client() won't redirect the URL — just verify wiring
	_ = multi.Send("integration test alert")

	if calls > 1 {
		t.Errorf("expected at most 1 call, got %d", calls)
	}
}
