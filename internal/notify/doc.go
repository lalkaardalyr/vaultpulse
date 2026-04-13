// Package notify provides clients for delivering alert messages to external
// notification services.
//
// Currently supported integrations:
//
//   - Slack — posts messages to a configured Incoming Webhook URL.
//
// Each client implements a simple Send(message string) error interface so that
// callers remain decoupled from the underlying transport.
package notify
