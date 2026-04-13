# vaultpulse

> A CLI tool that monitors HashiCorp Vault secret expiry and sends alerts before rotation deadlines.

---

## Installation

```bash
go install github.com/yourusername/vaultpulse@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultpulse.git
cd vaultpulse
go build -o vaultpulse .
```

---

## Usage

Set your Vault address and token, then run vaultpulse against a secret path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

# a secret path and alert if expiry is within 7 days
vaultpulse check --path secret/myapp/db-password --threshold 7d

# Watch multiple paths and send a Slack alert
vaultpulse watch \
  --path secret/myapp/api-key \
  --path secret/myapp/db-password \
  --threshold 14d \
  --alert slack \
  --webhook-url https://hooks.slack.com/services/xxx/yyy/zzz
```

### Available Commands

| Command   | Description                                      |
|-----------|--------------------------------------------------|
| `check`   | One-time check of a secret path                  |
| `watch`   | Continuously monitor paths on a schedule         |
| `list`    | List all monitored secrets and their expiry dates |

### Flags

```
--path          Vault secret path to monitor (repeatable)
--threshold     Alert window before expiry (e.g. 7d, 24h)
--alert         Alert backend: slack, email, or stdout (default: stdout)
--interval      Poll interval for watch mode (default: 1h)
```

---

## License

[MIT](LICENSE)