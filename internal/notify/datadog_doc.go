// Package notify provides notification sender implementations for VaultPulse.
//
// # Datadog
//
// The DatadogClient sends secret expiry alerts to the Datadog Events API.
// Alerts appear in the Datadog event stream tagged with "source:vaultpulse".
//
// Usage:
//
//	client, err := notify.NewDatadogClient(os.Getenv("DD_API_KEY"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("secret/db is expiring in 24 hours")
//
// The client sets the alert_type to "warning" and includes a standard title
// so events are easy to filter within Datadog dashboards and monitors.
package notify
