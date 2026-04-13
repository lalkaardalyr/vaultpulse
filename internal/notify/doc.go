// Package notify provides clients for sending alert notifications
// to external services such as Slack and PagerDuty.
//
// Each client implements the Sender interface defined in sender.go,
// allowing them to be used interchangeably within vaultpulse pipelines.
//
// Supported backends:
//   - Slack (incoming webhooks)
//   - PagerDuty (Events API v2)
package notify
