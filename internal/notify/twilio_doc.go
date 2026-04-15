// Package notify provides notification sender implementations for vaultpulse.
//
// # Twilio
//
// TwilioClient delivers secret-expiry alerts as SMS messages using the
// Twilio Messaging REST API (v2010-04-01).
//
// Required configuration:
//   - AccountSID  – Twilio account identifier (starts with "AC").
//   - AuthToken   – Secret credential paired with the AccountSID.
//   - From        – Twilio phone number or messaging service SID.
//   - To          – Destination phone number in E.164 format.
//
// Example:
//
//	client, err := notify.NewTwilioClient(sid, token, "+15550001111", "+15559998888")
//	if err != nil { ... }
//	err = client.Send("[vaultpulse] secret expiring in 24h")
package notify
