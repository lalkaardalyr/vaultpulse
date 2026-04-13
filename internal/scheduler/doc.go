// Package scheduler provides a periodic runner that coordinates secret
// monitoring and alert dispatch for vaultpulse.
//
// A Scheduler is constructed with a polling interval, a secrets.Monitor, and
// an alerts.Notifier. Calling Run blocks until the provided context is
// cancelled, executing a monitor pass on every tick and forwarding any
// resulting alerts to the notifier.
//
// Example usage:
//
//	s := scheduler.New(5*time.Minute, monitor, notifier)
//	if err := s.Run(ctx); err != nil && err != context.Canceled {
//		log.Fatalf("scheduler exited: %v", err)
//	}
package scheduler
