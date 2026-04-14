// Package main is the entry point for the vaultpulse CLI tool.
// It wires together configuration, scheduling, output formatting,
// audit logging, and multi-channel notifications into a single
// cohesive command-line application.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/vaultpulse/internal/audit"
	"github.com/yourusername/vaultpulse/internal/config"
	"github.com/yourusername/vaultpulse/internal/notify"
	"github.com/yourusername/vaultpulse/internal/runner"
	"github.com/yourusername/vaultpulse/internal/scheduler"
)

const version = "0
func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("vaultpulse", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		configPath  = fs.String("config", "config.yaml", "path to configuration file")
		showVersion = fs.Bool("version", false, "print version and exit")
		once        = fs.Bool("once", false, "run a single check and exit (ignores schedule interval)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Fprintf(stdout, "vaultpulse version %s\n", version)
		return nil
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Build optional notification senders based on config.
	var senders []notify.Sender

	if cfg.Notify.Slack.WebhookURL != "" {
		sc, err := notify.NewSlackClient(cfg.Notify.Slack.WebhookURL)
		if err != nil {
			return fmt.Errorf("slack client: %w", err)
		}
		senders = append(senders, sc)
	}

	if cfg.Notify.PagerDuty.RoutingKey != "" {
		pd, err := notify.NewPagerDutyClient(cfg.Notify.PagerDuty.RoutingKey)
		if err != nil {
			return fmt.Errorf("pagerduty client: %w", err)
		}
		senders = append(senders, pd)
	}

	if cfg.Notify.OpsGenie.APIKey != "" {
		og, err := notify.NewOpsGenieClient(cfg.Notify.OpsGenie.APIKey)
		if err != nil {
			return fmt.Errorf("opsgenie client: %w", err)
		}
		senders = append(senders, og)
	}

	if cfg.Notify.Webhook.URL != "" {
		wh, err := notify.NewWebhookClient(cfg.Notify.Webhook.URL)
		if err != nil {
			return fmt.Errorf("webhook client: %w", err)
		}
		senders = append(senders, wh)
	}

	// Set up audit logger (writes to stdout or a file if configured).
	auditWriter := stdout
	if cfg.Audit.LogFile != "" {
		f, err := os.OpenFile(cfg.Audit.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
		if err != nil {
			return fmt.Errorf("opening audit log: %w", err)
		}
		defer f.Close()
		auditWriter = f
	}
	auditLogger := audit.NewWithOptions(audit.WithWriter(auditWriter))
	_ = auditLogger // available for future hook-in to runner

	// Build the runner that performs a single scan cycle.
	r, err := runner.New(cfg, stdout, senders...)
	if err != nil {
		return fmt.Errorf("initialising runner: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if *once {
		return r.Run(ctx)
	}

	// Continuous mode: run on the configured interval.
	interval := time.Duration(cfg.Schedule.IntervalSeconds) * time.Second
	sched := scheduler.New(scheduler.WithInterval(interval))

	log.Printf("vaultpulse %s started — checking every %s", version, interval)

	return sched.Start(ctx, func(ctx context.Context) error {
		if err := r.Run(ctx); err != nil {
			// Log but do not stop the scheduler on transient errors.
			log.Printf("run error: %v", err)
		}
		return nil
	})
}
