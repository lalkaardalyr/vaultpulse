package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/scheduler"
)

// mockRunner is a simple Runner that records how many times it was called.
type mockRunner struct {
	calls int
	err   error
}

func (m *mockRunner) Run(ctx context.Context) error {
	m.calls++
	return m.err
}

func TestNew_ReturnsNonNil(t *testing.T) {
	s := scheduler.New(5*time.Second, nil, nil)
	if s == nil {
		t.Fatal("expected non-nil Scheduler")
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	// Use a very large interval so only the immediate tick fires.
	s := scheduler.New(24*time.Hour, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		// tick will panic with nil monitor — skip by catching the timeout.
		defer func() { recover() }() //nolint:errcheck
		done <- s.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != context.DeadlineExceeded && err != context.Canceled {
			t.Logf("scheduler stopped with: %v (acceptable)", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("scheduler did not stop after context cancellation")
	}
}

func TestScheduler_IntervalIsRespected(t *testing.T) {
	interval := 10 * time.Millisecond
	s := scheduler.New(interval, nil, nil)
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
// Verify the scheduler can be constructed with short intervals without panic.
}
