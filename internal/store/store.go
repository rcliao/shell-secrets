package store

// Store provides encrypted secret storage.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	List() ([]string, error)
	Remove(key string) error
	Close() error
}

// Option configures a FileStore.
type Option func(*FileStore)

// WithMasterKey provides an explicit master key instead of reading from the macOS Keychain.
func WithMasterKey(key []byte) Option {
	return func(fs *FileStore) {
		fs.masterKey = make([]byte, len(key))
		copy(fs.masterKey, key)
	}
}

// New creates a new FileStore. If storePath is empty, defaults to ~/.shell-secrets/secrets.enc.
func New(storePath string, opts ...Option) (Store, error) {
	return NewFileStore(storePath, opts...)
}
