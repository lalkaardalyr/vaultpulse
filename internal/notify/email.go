package notify

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// EmailClient sends alert notifications via SMTP.
type EmailClient struct {
	host     string
	port     int
	username stringword string
	from     string
	to       []string
}

// EmailConfig holds the configuration required to create an EmailClient.
type EmailConfig     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// NewEmailClient creates a new EmailClient from the provided config.
// Returns an error if required fields are missing.
func NewEmailClient(cfg EmailConfig) (*EmailClient, error) {
	if cfg.Host == "" {
		return nil, errors.New("email: host must be empty")
	}
	if cfg.From == "" {
		return nil, errors.New("email: from address must not be empty")
	}
	if len(cfg.To) == 0 {
		return nil, errors.New("email: at least one recipient is required")
	}
	port := cfg.Port
	if port == 0 {
		port = 587
	}
	return &EmailClient{
		host:     cfg.Host,
		port:     port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
		to:       cfg.To,
	}, nil
}

// Send delivers the alert message to all configured recipients via SMTP.
func (e *EmailClient) Send(message string) error {
	addr := fmt.Sprintf("%s:%d", e.host, e.port)

	var auth smtp.Auth
	if e.username != "" && e.password != "" {
		auth = smtp.PlainAuth("", e.username, e.password, e.host)
	}

	body := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: VaultPulse Alert\r\n\r\n%s",
		e.from,
		strings.Join(e.to, ", "),
		message,
	)

	if err := smtp.SendMail(addr, auth, e.from, e.to, []byte(body)); err != nil {
		return fmt.Errorf("email: failed to send message: %w", err)
	}
	return nil
}
