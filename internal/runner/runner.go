// Package runner wires together the core components of vaultpulse and
// executes a single monitoring cycle or a scheduled loop.
package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/example/vaultpulse/internal/alerts"
	"github.com/example/vaultpulse/internal/config"
	"github.com/example/vaultpulse/internal/output"
	"github.com/example/vaultpulse/internal/secrets"
	"github.com/example/vaultpulse/internal/vault"
)

// Runner orchestrates a single monitoring pass.
type Runner struct {
	cfg      *config.Config
	monitor  *secrets.Monitor
	notifier *alerts.Notifier
	formatter *output.Formatter
}

// New constructs a Runner from the provided config, writing output to w.
func New(cfg *config.Config, w io.Writer) (*Runner, error) {
	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("runner: create vault client: %w", err)
	}

	monitor := secrets.NewMonitor(client, cfg)
	notifier := alerts.NewNotifier(alerts.WithWriter(w))
	formatter := output.New(w)

	return &Runner{
		cfg:       cfg,
		monitor:   monitor,
		notifier:  notifier,
		formatter: formatter,
	}, nil
}

// Run executes one full monitoring cycle: fetch statuses, notify, and render.
func (r *Runner) Run(ctx context.Context) error {
	statuses, err := r.monitor.CheckAll(ctx)
	if err != nil {
		return fmt.Errorf("runner: check secrets: %w", err)
	}

	var alertList []alerts.Alert
	for _, s := range statuses {
		alertList = append(alertList, alerts.NewAlert(s))
	}

	if err := r.notifier.NotifyAll(alertList); err != nil {
		return fmt.Errorf("runner: notify: %w", err)
	}

	if r.cfg.Output.Format == "json" {
		return r.formatter.WriteJSON(statuses)
	}
	return r.formatter.WriteTable(statuses)
}
