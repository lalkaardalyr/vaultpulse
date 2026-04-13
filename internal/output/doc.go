// Package output provides formatting utilities for rendering secret status
// reports to the terminal or other writers.
//
// Two formats are supported:
//
//   - table: a human-readable tabular layout printed to stdout (default).
//   - json: a machine-readable JSON array suitable for piping to other tools.
//
// Usage:
//
//	f := output.New(output.FormatTable)
//	f.Write(statuses)
package output
