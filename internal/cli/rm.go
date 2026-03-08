package cli

import (
	"fmt"

	"github.com/rcliao/shell-secrets/internal/store"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <key>",
	Short: "Remove a secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New(storePath)
		if err != nil {
			return err
		}
		defer s.Close()

		if err := s.Remove(args[0]); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Secret %q removed.\n", args[0])
		return nil
	},
}
