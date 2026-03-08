package cli

import (
	"fmt"

	"github.com/rcliao/shell-secrets/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all secret key names",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New(storePath)
		if err != nil {
			return err
		}
		defer s.Close()

		keys, err := s.List()
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No secrets stored.")
			return nil
		}

		for _, k := range keys {
			fmt.Fprintln(cmd.OutOrStdout(), k)
		}
		return nil
	},
}
