package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rcliao/shell-secrets/internal/crypto"
)

func testStore(t *testing.T) *FileStore {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "secrets.enc")
	fs, err := NewFileStore(path, WithMasterKey(key))
	if err != nil {
		t.Fatal(err)
	}
	return fs
}

func TestSetAndGet(t *testing.T) {
	fs := testStore(t)

	if err := fs.Set("API_KEY", "secret123"); err != nil {
		t.Fatal(err)
	}

	val, err := fs.Get("API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if val != "secret123" {
		t.Fatalf("expected %q, got %q", "secret123", val)
	}
}

func TestGetMissing(t *testing.T) {
	fs := testStore(t)

	_, err := fs.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestList(t *testing.T) {
	fs := testStore(t)

	fs.Set("B_KEY", "val")
	fs.Set("A_KEY", "val")

	keys, err := fs.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "B_KEY" {
		t.Fatalf("expected sorted keys [A_KEY B_KEY], got %v", keys)
	}
}

func TestRemove(t *testing.T) {
	fs := testStore(t)

	fs.Set("KEY", "val")
	if err := fs.Remove("KEY"); err != nil {
		t.Fatal(err)
	}

	_, err := fs.Get("KEY")
	if err == nil {
		t.Fatal("expected error after remove")
	}
}

func TestRemoveMissing(t *testing.T) {
	fs := testStore(t)

	err := fs.Remove("MISSING")
	if err == nil {
		t.Fatal("expected error removing missing key")
	}
}

func TestPersistence(t *testing.T) {
	key, _ := crypto.GenerateKey()
	path := filepath.Join(t.TempDir(), "secrets.enc")

	// Create and write
	fs1, err := NewFileStore(path, WithMasterKey(key))
	if err != nil {
		t.Fatal(err)
	}
	fs1.Set("TOKEN", "abc123")
	fs1.Close()

	// Reopen and read
	fs2, err := NewFileStore(path, WithMasterKey(key))
	if err != nil {
		t.Fatal(err)
	}
	val, err := fs2.Get("TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if val != "abc123" {
		t.Fatalf("expected %q, got %q", "abc123", val)
	}
}

func TestFilePermissions(t *testing.T) {
	fs := testStore(t)
	fs.Set("KEY", "val")

	info, err := os.Stat(fs.path)
	if err != nil {
		t.Fatal(err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Fatalf("expected file permissions 0600, got %o", perm)
	}
}

func TestOverwrite(t *testing.T) {
	fs := testStore(t)

	fs.Set("KEY", "first")
	fs.Set("KEY", "second")

	val, _ := fs.Get("KEY")
	if val != "second" {
		t.Fatalf("expected %q, got %q", "second", val)
	}
}
