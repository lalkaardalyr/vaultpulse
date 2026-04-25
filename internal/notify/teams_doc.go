// Package notify provides notification clients for delivering VaultPulse
// alerts to external services.
//
// # Microsoft Teams
//
// TeamsClient sends alert messages to a Microsoft Teams channel using an
// incoming webhook URL configured in the Teams channel settings.
//
// The webhook URL must be a valid HTTPS URL obtained from the Teams channel
// connector configuration. Messages are delivered as plain text cards.
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
//
// # Notification Severity Levels
//
// Alert messages should be prefixed with a severity tag to aid visibility
// in the Teams channel:
//
//	[CRITICAL] - secret expires within 1 day
//	[WARNING]  - secret expires within 7 days
//	[INFO]     - secret expires within 30 days
package notify
