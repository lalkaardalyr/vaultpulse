package notify

import (
	"errors"
	"testing"
)

// fakeSender is a lightweight Sender stub for integration-style tests.
type fakeSender struct {
	sentMsg string
	retErr  error
}

func (f *fakeSender) Send(msg string) error {
	f.sentMsg = msg
	return f.retErr
}

func TestMultiSender_WithSNSLike_SendsAll(t *testing.T) {
	a := &fakeSender{}
	b := &fakeSender{}

	multi, err := NewMultiSender(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const msg = "vault secret expiring in 24h"
	if err := multi.Send(msg); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
	if a.sentMsg != msg {
		t.Errorf("sender A: got %q, want %q", a.sentMsg, msg)
	}
	if b.sentMsg != msg {
		t.Errorf("sender B: got %q, want %q", b.sentMsg, msg)
	}
}

func TestMultiSender_WithSNSLike_AccumulatesErrors(t *testing.T) {
	a := &fakeSender{retErr: errors.New("sns unavailable")}
	b := &fakeSender{}

	multi, err := NewMultiSender(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := multi.Send("alert"); err == nil {
		t.Fatal("expected aggregated error")
	}
	// sender b should still have been called
	if b.sentMsg != "alert" {
		t.Errorf("expected sender B to be called despite sender A error")
	}
}
