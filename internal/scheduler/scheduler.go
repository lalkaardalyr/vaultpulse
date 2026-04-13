package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/vaultpulse/internal/alerts"
	"github.com/yourusername/vaultpulse/internal/secrets"
)

// Runner defines the interface for running a scheduled monitor check.
type Runner interface {
	Run(ctx context.Context) error
}

// Scheduler periodically triggers secret monitoring and alerting.
type Scheduler struct {
	interval time.Duration
	monitor  *secrets.Monitor
	notifier *alerts.Notifier
}

// New creates a new Scheduler with the given interval, monitor, and notifier.
func New(interval time.Duration, monitor *secrets.Monitor, notifier *alerts.Notifier) *Scheduler {
	return &Scheduler{
		interval: interval,
		monitor:  monitor,
		notifier: notifier,
	}
}

// Run starts the scheduling loop and blocks until the context is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	log.Printf("scheduler: starting with interval %s", s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run once immediately before waiting for the first tick.
	if err := s.tick(ctx); err != nil {
		log.Printf("scheduler: tick error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := s.tick(ctx); err != nil {
				log.Printf("scheduler: tick error: %v", err)
			}
		case <-ctx.Done():
			log.Println("scheduler: context cancelled, stopping")
			return ctx.Err()
		}
	}
}

// tick performs a single monitor pass and dispatches alerts.
func (s *Scheduler) tick(ctx context.Context) error {
	statuses, err := s.monitor.CheckAll(ctx)
	if err != nil {
		return err
	}
	s.notifier.NotifyAll(statuses)
	return nil
}
