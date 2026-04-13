// Package alerts provides alert construction and notification delivery
// for VaultPulse secret expiry monitoring.
//
// Alerts are classified into three levels based on configurable time
// thresholds relative to a secret's expiry:
//
//   - INFO     – expiry is beyond the warning threshold (no action needed)
//   - WARNING  – expiry is within the warning threshold
//   - CRITICAL – expiry is within the critical threshold
//
// Usage:
//
//	notifier := alerts.NewNotifier(48*time.Hour, 24*time.Hour)
//	notifier.NotifyAll(secretExpiryMap)
package alerts
