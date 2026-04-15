package notify

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// TestSignalSciencesClient_MultiSender_Integration verifies that
// SignalSciencesClient works correctly when composed inside a MultiSender.
func TestSignalSciencesClient_MultiSender_Integration(t *testing.T) {
	var callCount int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c1, err := NewSignalSciencesClient("user@example.com", "token1", "corp-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c1.endpoint = ts.URL

	c2, err := NewSignalSciencesClient("user@example.com", "token2", "corp-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c2.endpoint = ts.URL

	multi, err := NewMultiSender(c1, c2)
	if err != nil {
		t.Fatalf("unexpected error creating multi sender: %v", err)
	}

	if err := multi.Send("integration test alert"); err != nil {
		t.Fatalf("unexpected error from multi sender: %v", err)
	}

	if got := atomic.LoadInt32(&callCount); got != 2 {
		t.Errorf("expected 2 calls to server, got %d", got)
	}
}
