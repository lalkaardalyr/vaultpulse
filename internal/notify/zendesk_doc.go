// Package notify provides notification sender implementations for VaultPulse.
//
// # Zendesk
//
// ZendeskClient creates a high-priority support ticket in a Zendesk account
// whenever a VaultPulse alert is triggered.
//
// Usage:
//
//	client, err := notify.NewZendeskClient("mycompany", "user@example.com", "apitoken")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Secret at secret/db is expiring in 2 days")
package notify
