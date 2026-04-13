// Package runner wires together the Vault client, secrets monitor, alert
// notifier, and output formatter into a single monitoring pass.
//
// Typical usage:
//
//	r, err := runner.New(cfg, os.Stdout)
//	if err != nil { ... }
//	if err := r.Run(ctx); err != nil { ... }
//
// The runner is designed to be invoked either once (one-shot mode) or
// repeatedly by the scheduler for continuous monitoring.
package runner
