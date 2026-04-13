package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration for VaultPulse.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Monitor MonitorConfig `yaml:"monitor"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	TLSSkip   bool   `yaml:"tls_skip_verify"`
}

// AlertsConfig defines alerting thresholds and channels.
type AlertsConfig struct {
	WarningDays  int    `yaml:"warning_days"`
	CriticalDays int    `yaml:"critical_days"`
	SlackWebhook string `yaml:"slack_webhook"`
	EmailSMTP    string `yaml:"email_smtp"`
	EmailTo      string `yaml:"email_to"`
}

// MonitorConfig defines which secret paths to monitor.
type MonitorConfig struct {
	Paths    []string      `yaml:"paths"`
	Interval time.Duration `yaml:"interval"`
}

// Load reads and parses the config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" {
		c.Vault.Token = os.Getenv("VAULT_TOKEN")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("vault.token or VAULT_TOKEN env var is required")
	}
	if len(c.Monitor.Paths) == 0 {
		return fmt.Errorf("monitor.paths must contain at least one path")
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Alerts.WarningDays == 0 {
		c.Alerts.WarningDays = 14
	}
	if c.Alerts.CriticalDays == 0 {
		c.Alerts.CriticalDays = 3
	}
	if c.Monitor.Interval == 0 {
		c.Monitor.Interval = 1 * time.Hour
	}
}
