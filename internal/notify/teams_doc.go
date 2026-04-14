// Package notify provides notification clients for delivering VaultPulse
// alerts to external services.
//
// # Microsoft Teams
//
// TeamsClient sends alert messages to a Microsoft Teams channel using an
// incoming webhook URL configured in the Teams channel settings.
//
// Example usage:
//
//	client, err := notify.NewTeamsClient("https://outlook.office.com/webhook/...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("[CRITICAL] secret/db/prod expires in 1 day"); err != nil {
//		log.Printf("teams notification failed: %v", err)
//	}
package notify
