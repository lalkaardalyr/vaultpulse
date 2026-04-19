// Package notify provides notification sender implementations for vaultpulse.
//
// # Sentry
//
// SentryClient sends alert messages to a Sentry DSN endpoint using the
// Sentry HTTP API. It is suitable for capturing secret expiry events as
// Sentry error-level messages for visibility in error tracking dashboards.
//
// Usage:
//
//	client, err := notify.NewSentryClient("https://...@sentry.io/...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("secret expiring soon: secret/myapp/db")
package notify
