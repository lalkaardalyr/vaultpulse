// Package notify provides notification sender implementations for VaultPulse alerts.
//
// # Signal Sciences
//
// SignalSciencesClient sends alert messages to the Signal Sciences (Fastly Next-Gen WAF)
// custom event ingestion API. Each alert is posted as a custom event tagged with
// "vaultpulse" and "secret-expiry" for easy filtering in the dashboard.
//
// Required configuration:
//   - APIUser:  Signal Sciences API user (email address)
//   - APIToken: Signal Sciences API token
//   - CorpName: The corp slug identifying your organisation
//
// Usage:
//
//	client, err := notify.NewSignalSciencesClient(apiUser, apiToken, corpName)
//	if err != nil { ... }
//	err = client.Send("secret expiring soon: secret/my-app/db")
package notify
