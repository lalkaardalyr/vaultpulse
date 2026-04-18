// Package notify provides notification clients for various alerting platforms.
//
// SIGNL4 (signl4.go) integrates with the SIGNL4 mobile alerting service.
// It sends structured JSON payloads to a SIGNL4 team webhook URL.
//
// Usage:
//
//	client, err := notify.NewSIGNL4Client("https://connect.signl4.com/webhook/<team-secret>")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Secret expiry warning: secret/myapp/db expires in 24h")
package notify
