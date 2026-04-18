// Package notify provides notification clients for various alerting platforms.
//
// # Statuspage
//
// StatuspageClient sends incident alerts to Atlassian Statuspage via its REST API.
// It creates a new incident with status "investigating" whenever a VaultPulse
// alert is triggered.
//
// Usage:
//
//	client, err := notify.NewStatuspageClient(apiKey, pageID)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Secret at secret/db/password expires in 24h")
package notify
