package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatrixClient_MultiSender_Integration(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$ev"}`))
	}))
	defer server.Close()

	matrix, err := NewMatrixClient(server.URL, "tok", "!r:matrix.org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	multi, err := NewMultiSender(matrix)
	if err != nil {
		t.Fatalf("unexpected error building MultiSender: %v", err)
	}

	if err := multi.Send("integration test"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call to Matrix server, got %d", calls)
	}
}
