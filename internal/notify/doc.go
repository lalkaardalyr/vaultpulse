// Package notify provides notification backends for VaultPulse alerts.
//
// Supported backends:
//   - Slack (via incoming webhook URL)
//   - PagerDuty (via Events API v2)
//   - Email (via SMTP)
//   - Generic HTTP Webhook
//
// All backends implement the Sender interface, allowing them to be used
// interchangeably by the alert pipeline.
package notify
