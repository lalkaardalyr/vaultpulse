package notify

import (
	"errors"
	"testing"
)

// stubSender is a test double for Sender.
type stubSender struct {
	called  bool
	recv    string
	errOnce error
}

func (s *stubSender) Send(msg string) error {
	s.called = true
	s.recv = msg
	return s.errOnce
}

func TestNewMultiSender_NoSenders_ReturnsError(t *testing.T) {
	_, err := NewMultiSender()
	if err == nil {
		t.Fatal("expected error when no senders provided")
	}
}

func TestNewMultiSender_ValidSenders_ReturnsMulti(t *testing.T) {
	m, err := NewMultiSender(&stubSender{}, &stubSender{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil MultiSender")
	}
}

func TestMultiSender_Send_CallsAllSenders(t *testing.T) {
	a, b := &stubSender{}, &stubSender{}
	m, _ := NewMultiSender(a, b)

	if err := m.Send("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both senders to be called")
	}
	if a.recv != "hello" || b.recv != "hello" {
		t.Error("expected both senders to receive the message")
	}
}

func TestMultiSender_Send_AccumulatesErrors(t *testing.T) {
	a := &stubSender{errOnce: errors.New("sender a failed")}
	b := &stubSender{errOnce: errors.New("sender b failed")}
	m, _ := NewMultiSender(a, b)

	err := m.Send("msg")
	if err == nil {
		t.Fatal("expected combined error, got nil")
	}
	if !errors.Is(err, a.errOnce) || !errors.Is(err, b.errOnce) {
		t.Errorf("expected both errors in result, got: %v", err)
	}
}

func TestMultiSender_Send_PartialError_StillCallsAll(t *testing.T) {
	a := &stubSender{errOnce: errors.New("oops")}
	b := &stubSender{}
	m, _ := NewMultiSender(a, b)

	_ = m.Send("ping")
	if !b.called {
		t.Error("expected second sender to be called even after first fails")
	}
}
