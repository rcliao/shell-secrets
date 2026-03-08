package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rcliao/shell-secrets/internal/store"
	"github.com/spf13/cobra"
)

var useStdin bool

var setCmd = &cobra.Command{
	Use:   "set <key> [value]",
	Short: "Store a secret",
	Long:  "Store a secret. Use --stdin to read the value from stdin (avoids shell history).",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		var value string

		if useStdin {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				value = scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
		} else if len(args) == 2 {
			value = args[1]
		} else {
			return fmt.Errorf("provide a value as argument or use --stdin")
		}

		value = strings.TrimSpace(value)
		if value == "" {
			return fmt.Errorf("value cannot be empty")
		}

		s, err := store.New(storePath)
		if err != nil {
			return err
		}
		defer s.Close()

		if err := s.Set(key, value); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Secret %q stored.\n", key)
		return nil
	},
}

func init() {
	setCmd.Flags().BoolVar(&useStdin, "stdin", false, "read value from stdin")
}
