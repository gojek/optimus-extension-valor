package cmd

import (
	"os"

	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/spf13/cobra"
)

const (
	defaultRecipeType   = "file"
	defaultRecipeFormat = "yaml"
	defaultRecipePath   = "./valor.yaml"

	defaultBatchSize = 4
)

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

func enrichWithBatchSize(r *recipe.Recipe) error {
	for i := 0; i < len(r.Resources); i++ {
		if r.Resources[i].BatchSize == 0 {
			r.Resources[i].BatchSize = defaultBatchSize
		}
	}
	return nil
}
