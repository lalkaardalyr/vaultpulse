// Package notify provides integrations with external notification services.
//
// ClickSend
//
// ClickSendClient delivers SMS alerts through the ClickSend REST API.
// Credentials (username and API key) are required along with a destination
// phone number in E.164 format (e.g. "+14155552671").
//
// Example:
//
//	client, err := notify.NewClickSendClient("user@example.com", "API_KEY", "+14155552671")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Vault secret expiring in 24 hours")
package notify
