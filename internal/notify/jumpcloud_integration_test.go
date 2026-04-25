package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestJumpCloudClient_MultiSender_Integration(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	jc, err := notify.NewJumpCloudClient("key-abc", "org-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	jc.SetEndpoint(ts.URL)

	slack, err := notify.NewSlackClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	multi, err := notify.NewMultiSender(jc, slack)
	if err != nil {
		t.Fatalf("unexpected error creating multi sender: %v", err)
	}

	if err := multi.Send("integration test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calls != 2 {
		t.Errorf("expected 2 HTTP calls, got %d", calls)
	}
}
