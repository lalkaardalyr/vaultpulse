// Package audit implements structured, newline-delimited JSON audit logging
// for VaultPulse. Every significant event — scan lifecycle, secret expiry
// detection, and alert dispatch — is recorded with a timestamp, event type,
// optional secret path, human-readable message, and arbitrary metadata.
//
// Usage:
//
//	logger := audit.New(os.Stderr)
//	logger.Log(audit.EventScanStarted, "", "scan initiated", nil)
package audit
