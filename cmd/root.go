package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	defaultRecipeType   = "file"
	defaultRecipeFormat = "yaml"
)

const defaultRecipePath = "./valor.yaml"

var recipePath string

// Execute executes command
func Execute() {
	rootCmd := &cobra.Command{
		Use:          "valor",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(getExecuteCmd())
	rootCmd.AddCommand(getProfileCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
