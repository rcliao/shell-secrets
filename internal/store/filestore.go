package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/rcliao/shell-secrets/internal/crypto"
	"github.com/rcliao/shell-secrets/internal/keychain"
)

type encryptedFile struct {
	Version    int    `json:"version"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

// FileStore implements Store with AES-256-GCM encryption backed by a JSON file.
type FileStore struct {
	mu        sync.RWMutex
	path      string
	masterKey []byte
	data      map[string]string
}

// NewFileStore creates a FileStore. If storePath is empty, defaults to ~/.shell-secrets/secrets.enc.
func NewFileStore(storePath string, opts ...Option) (*FileStore, error) {
	if storePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home dir: %w", err)
		}
		storePath = filepath.Join(home, ".shell-secrets", "secrets.enc")
	}

	fs := &FileStore{
		path: storePath,
		data: make(map[string]string),
	}

	for _, opt := range opts {
		opt(fs)
	}

	// Load master key from Keychain if not provided
	if fs.masterKey == nil {
		key, err := keychain.Load()
		if err != nil {
			return nil, err
		}
		fs.masterKey = key
	}

	// Load existing data if file exists
	if err := fs.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStore) Get(key string) (string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	val, ok := fs.data[key]
	if !ok {
		return "", fmt.Errorf("secret %q not found", key)
	}
	return val, nil
}

func (fs *FileStore) Set(key, value string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.data[key] = value
	return fs.save()
}

func (fs *FileStore) List() ([]string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	keys := make([]string, 0, len(fs.data))
	for k := range fs.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func (fs *FileStore) Remove(key string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, ok := fs.data[key]; !ok {
		return fmt.Errorf("secret %q not found", key)
	}
	delete(fs.data, key)
	return fs.save()
}

func (fs *FileStore) Close() error {
	return nil
}

func (fs *FileStore) load() error {
	raw, err := os.ReadFile(fs.path)
	if err != nil {
		return err
	}

	var ef encryptedFile
	if err := json.Unmarshal(raw, &ef); err != nil {
		return fmt.Errorf("parsing encrypted file: %w", err)
	}

	if ef.Version != 1 {
		return fmt.Errorf("unsupported file version: %d", ef.Version)
	}

	nonce, err := base64.StdEncoding.DecodeString(ef.Nonce)
	if err != nil {
		return fmt.Errorf("decoding nonce: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ef.Ciphertext)
	if err != nil {
		return fmt.Errorf("decoding ciphertext: %w", err)
	}

	plaintext, err := crypto.Decrypt(fs.masterKey, nonce, ciphertext)
	if err != nil {
		return fmt.Errorf("decrypting secrets: %w", err)
	}

	if err := json.Unmarshal(plaintext, &fs.data); err != nil {
		return fmt.Errorf("parsing decrypted secrets: %w", err)
	}

	return nil
}

func (fs *FileStore) save() error {
	plaintext, err := json.Marshal(fs.data)
	if err != nil {
		return fmt.Errorf("marshaling secrets: %w", err)
	}

	nonce, ciphertext, err := crypto.Encrypt(fs.masterKey, plaintext)
	if err != nil {
		return err
	}

	ef := encryptedFile{
		Version:    1,
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}

	raw, err := json.MarshalIndent(ef, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling encrypted file: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fs.path), 0700); err != nil {
		return fmt.Errorf("creating store directory: %w", err)
	}

	if err := os.WriteFile(fs.path, raw, 0600); err != nil {
		return fmt.Errorf("writing encrypted file: %w", err)
	}

	return nil
}
