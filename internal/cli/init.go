package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rcliao/shell-secrets/internal/crypto"
	"github.com/rcliao/shell-secrets/internal/keychain"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate master key and store in macOS Keychain",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Generate master key
		key, err := crypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("generating master key: %w", err)
		}

		// Store in Keychain
		if err := keychain.Store(key); err != nil {
			return fmt.Errorf("storing master key: %w", err)
		}

		// Create store directory
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home dir: %w", err)
		}
		storeDir := filepath.Join(home, ".shell-secrets")
		if err := os.MkdirAll(storeDir, 0700); err != nil {
			return fmt.Errorf("creating store directory: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Master key generated and stored in macOS Keychain.")
		fmt.Fprintf(cmd.OutOrStdout(), "Store directory: %s\n", storeDir)
		return nil
	},
}
