# shell-secrets

Lightweight encrypted secret store. AES-256-GCM encryption with macOS Keychain integration.

Part of the [Ghost in the Shell](https://github.com/rcliao?tab=repositories&q=shell-) ecosystem.

## Install

```bash
go install github.com/rcliao/shell-secrets/cmd/shell-secrets@latest
```

## Quick Start

```bash
# Initialize (generates master key, stores in macOS Keychain)
shell-secrets init

# Store a secret
shell-secrets set API_KEY "sk-abc123"

# Store from stdin (avoids shell history)
echo "sk-abc123" | shell-secrets set API_KEY --stdin

# Retrieve a secret
shell-secrets get API_KEY

# List all secret names
shell-secrets list

# Remove a secret
shell-secrets rm API_KEY
```

## How It Works

1. `shell-secrets init` generates a random 32-byte AES-256 key and stores it in the macOS Keychain
2. Secrets are encrypted with AES-256-GCM and stored in `~/.shell-secrets/secrets.enc`
3. Each operation loads the master key from Keychain, encrypts/decrypts as needed

## Security

- AES-256-GCM authenticated encryption
- Random nonce per encryption
- Master key stored in macOS Keychain (not on disk)
- File permissions `0600` (owner-only)

## Library Usage

```go
import secrets "github.com/rcliao/shell-secrets"

store, err := secrets.NewStore("")  // default path, key from Keychain
val, err := store.Get("API_KEY")
err = store.Set("API_KEY", "value")
```

## Build

```bash
make build    # Build binary
make test     # Run tests
make vet      # Run go vet
```

## License

MIT
