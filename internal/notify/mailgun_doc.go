// Package notify provides alert sender implementations for VaultPulse.
//
// # Mailgun
//
// MailgunClient delivers alert messages via the Mailgun transactional email API.
// It requires a Mailgun domain, API key, a sender address, and a recipient address.
//
// Example usage:
//
//	client, err := notify.NewMailgunClient(
//		"mg.example.com",
//		"key-xxxxxxxxxxxxxxxx",
//		"alerts@mg.example.com",
//		"oncall@example.com",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	client.Send("Secret expiry warning: secret/db/password expires in 3 days")
package notify
