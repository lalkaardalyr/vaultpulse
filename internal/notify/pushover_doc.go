// Package notify provides integrations with external notification services
// for delivering VaultPulse secret-expiry alerts.
//
// # Pushover
//
// PushoverClient sends alert messages to mobile and desktop devices via the
// Pushover (https://pushover.net) API. Each message is tagged with the
// "VaultPulse Alert" title so recipients can filter notifications easily.
//
// Usage:
//
//	client, err := notify.NewPushoverClient(apiToken, userKey)
//	if err != nil { ... }
//	err = client.Send("secret/my-app expires in 3 days")
package notify
