package cli

import (
	"github.com/spf13/cobra"
)

var storePath string

var RootCmd = &cobra.Command{
	Use:   "shell-secrets",
	Short: "Encrypted secret manager",
}

func init() {
	RootCmd.PersistentFlags().StringVar(&storePath, "store-path", "", "path to encrypted store file (default ~/.shell-secrets/secrets.enc)")

	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(setCmd)
	RootCmd.AddCommand(getCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(rmCmd)
}
