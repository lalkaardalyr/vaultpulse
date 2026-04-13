// Package notify provides notification clients for delivering VaultPulse
// alerts to external services.
//
// EmailClient sends alerts via SMTP. It supports plain-text authentication
// and defaults to port 587 when no port is specified. Configure it using
// EmailConfig and create an instance with NewEmailClient.
//
// Example usage:
//
//	client, err := notify.NewEmailClient(notify.EmailConfig{
//		Host:     "smtp.example.com",
//		Port:     587,
//		Username: "user@example.com",
//		Password: "secret",
//		From:     "vaultpulse@example.com",
//		To:       []string{"ops@example.com"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("Secret expiry warning"); err != nil {
//		log.Println("alert delivery failed:", err)
//	}
package notify
