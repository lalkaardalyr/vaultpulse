// Package notify provides notification sender implementations for VaultPulse.
//
// # New Relic
//
// NewRelicClient sends alert events to New Relic Insights using the
// custom events API. Each alert is recorded as a "VaultPulseAlert" event
// type, making it queryable via NRQL.
//
// Usage:
//
//	client, err := notify.NewNewRelicClient(accountID, apiKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("secret expiry warning: secret/db/prod expires in 2 days")
package notify
