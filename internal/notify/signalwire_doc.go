// Package notify provides integrations for sending alert notifications
// to various third-party services.
//
// SignalWireClient sends SMS alerts via the SignalWire messaging API.
// It requires a project ID, auth token, space URL, sender number, and
// recipient number.
//
// Usage:
//
//	client, err := notify.NewSignalWireClient(
//		"project-id",
//		"auth-token",
//		"example.signalwire.com",
//		"+15550001111",
//		"+15559998888",
//	)
package notify
