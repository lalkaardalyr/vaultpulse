package output

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/vaultpulse/internal/secrets"
)

// Format controls the output format of secret status reports.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Formatter writes secret status output to a writer.
type Formatter struct {
	w      io.Writer
	format Format
}

// New creates a Formatter with the given format, defaulting to stdout.
func New(format Format) *Formatter {
	return &Formatter{w: os.Stdout, format: format}
}

// WithWriter returns a copy of the Formatter writing to w.
func (f *Formatter) WithWriter(w io.Writer) *Formatter {
	return &Formatter{w: w, format: f.format}
}

// Write outputs the list of SecretStatus entries using the configured format.
func (f *Formatter) Write(statuses []secrets.SecretStatus) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(statuses)
	default:
		return f.writeTable(statuses)
	}
}

func (f *Formatter) writeTable(statuses []secrets.SecretStatus) error {
	tw := tabwriter.NewWriter(f.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tSTATUS\tEXPIRES IN\tEXPIRY")
	fmt.Fprintln(tw, "----\t------\t----------\t------")
	for _, s := range statuses {
		expiry := s.ExpiresAt.Format(time.RFC3339)
		ttl := time.Until(s.ExpiresAt).Round(time.Second)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.Path, s.Status, ttl, expiry)
	}
	return tw.Flush()
}

func (f *Formatter) writeJSON(statuses []secrets.SecretStatus) error {
	fmt.Fprintln(f.w, "[")
	for i, s := range statuses {
		comma := ","
		if i == len(statuses)-1 {
			comma = ""
		}
		fmt.Fprintf(f.w, "  {\"path\":%q,\"status\":%q,\"expires_at\":%q}%s\n",
			s.Path, s.Status, s.ExpiresAt.Format(time.RFC3339), comma)
	}
	fmt.Fprintln(f.w, "]")
	return nil
}
