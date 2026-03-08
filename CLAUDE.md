# shell-secrets

Lightweight encrypted secret store. AES-256-GCM with macOS Keychain.

## Architecture

- `cmd/shell-secrets/main.go` — Cobra CLI entrypoint
- `internal/cli/` — Commands: init, set, get, list, rm
- `internal/store/` — FileStore: encrypted JSON persistence at `~/.shell-secrets/secrets.enc`
- `internal/crypto/` — AES-256-GCM encrypt/decrypt with random nonce
- `internal/keychain/` — macOS Keychain integration (service: `shell-secrets`)
- `secrets.go` — Public API: `Store`, `NewStore()`, `WithMasterKey()`

## Build & Test

```bash
make build    # Build binary
make test     # Run tests
make vet      # Run go vet
```

## Key Patterns

- Master key stored in macOS Keychain, never on disk
- File format: JSON with version, nonce, base64 ciphertext
- `--stdin` flag for set command avoids shell history exposure
- `--store-path` flag overrides default `~/.shell-secrets/secrets.enc`
- `WithMasterKey()` option for testing without Keychain
- Thread-safe with `sync.RWMutex`
