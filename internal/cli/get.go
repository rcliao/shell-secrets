package cli

import (
	"fmt"

	"github.com/rcliao/shell-secrets/internal/store"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Print a decrypted secret to stdout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New(storePath)
		if err != nil {
			return err
		}
		defer s.Close()

		val, err := s.Get(args[0])
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), val)
		return nil
	},
}
