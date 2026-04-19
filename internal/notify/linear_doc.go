// Package notify provides notification clients for VaultPulse alerts.
//
// LinearClient sends alert notifications to Linear by creating issues
// via the Linear GraphQL API. Each alert message becomes the title of
// a new issue assigned to the configured team.
//
// Usage:
//
//	client, err := notify.NewLinearClient(apiKey, teamID)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("secret expiring in 24h: secret/db/prod")
package notify
