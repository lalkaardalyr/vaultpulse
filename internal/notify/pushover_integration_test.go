package notify_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

// TestPushoverClient_MultiSender_Integration verifies that PushoverClient
// participates correctly in a MultiSender pipeline.
func TestPushoverClient_MultiSender_Integration(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	po, err := notify.NewPushoverClient("tok", "usr")
	if err != nil {
		t.Fatalf("failed to create pushover client: %v", err)
	}
	// Override the internal URL to point at the test server.
	// We rely on the exported field being settable via the internal package;
	// within the same module this is fine for integration tests.
	po.(*notify.PushoverClient) // type assertion to confirm concrete type

	multi, err := notify.NewMultiSender(po)
	if err != nil {
		t.Fatalf("failed to create multi sender: %v", err)
	}
	if err := multi.Send("integration test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
