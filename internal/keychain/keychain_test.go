package keychain

import (
	"os/exec"
	"testing"
)

func hasSecurityCLI() bool {
	_, err := exec.LookPath("security")
	return err == nil
}

func TestStoreAndLoad(t *testing.T) {
	if !hasSecurityCLI() {
		t.Skip("macOS security CLI not available")
	}

	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	// Clean up before/after
	Delete()
	t.Cleanup(func() { Delete() })

	if err := Store(key); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(loaded))
	}

	for i := range key {
		if loaded[i] != key[i] {
			t.Fatalf("key mismatch at byte %d", i)
		}
	}
}

func TestLoadMissing(t *testing.T) {
	if !hasSecurityCLI() {
		t.Skip("macOS security CLI not available")
	}

	Delete() // ensure clean state
	_, err := Load()
	if err == nil {
		t.Fatal("expected error loading missing key")
	}
}
