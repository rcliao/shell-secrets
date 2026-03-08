package keychain

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

const (
	serviceName = "shell-secrets"
	accountName = "master-key"
)

// Store saves the master key to the macOS Keychain.
func Store(key []byte) error {
	encoded := hex.EncodeToString(key)

	// Delete any existing entry first (ignore errors).
	exec.Command("security", "delete-generic-password",
		"-s", serviceName, "-a", accountName).Run()

	cmd := exec.Command("security", "add-generic-password",
		"-s", serviceName, "-a", accountName,
		"-w", encoded,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("keychain store: %w: %s", err, out)
	}
	return nil
}

// Load retrieves the master key from the macOS Keychain.
func Load() ([]byte, error) {
	cmd := exec.Command("security", "find-generic-password",
		"-s", serviceName, "-a", accountName,
		"-w", // Output password only
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("keychain load: master key not found (run 'shell-secrets init' first)")
	}

	encoded := strings.TrimSpace(string(out))
	key, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("keychain load: invalid key format: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("keychain load: expected 32-byte key, got %d", len(key))
	}

	return key, nil
}

// Delete removes the master key from the macOS Keychain.
func Delete() error {
	cmd := exec.Command("security", "delete-generic-password",
		"-s", serviceName, "-a", accountName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("keychain delete: %w: %s", err, out)
	}
	return nil
}
