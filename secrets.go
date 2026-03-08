package secrets

import (
	"github.com/rcliao/shell-secrets/internal/store"
)

// Store provides encrypted secret storage.
type Store = store.Store

// Option configures a Store.
type Option = store.Option

// NewStore creates a new encrypted secret store at the given path.
// If storePath is empty, defaults to ~/.shell-secrets/secrets.enc.
func NewStore(storePath string, opts ...Option) (Store, error) {
	return store.New(storePath, opts...)
}

// WithMasterKey provides an explicit master key instead of reading from the macOS Keychain.
func WithMasterKey(key []byte) Option {
	return store.WithMasterKey(key)
}
