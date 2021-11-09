package main

import (
	"errors"
	"os"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/model"
	_ "github.com/gojek/optimus-extension-valor/plugin/endec"
	_ "github.com/gojek/optimus-extension-valor/plugin/formatter"
	_ "github.com/gojek/optimus-extension-valor/plugin/io"
	_ "github.com/gojek/optimus-extension-valor/plugin/progress"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/endec"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const (
	defaultRecipeType   = "file"
	defaultRecipePath   = "./valor.yaml"
	defaultRecipeFormat = "yaml"
	defaultBatchSize    = 4
	defaultProgressType = "verbose"
)

func main() {
	var path string
	var batch int
	var progressType string
	cmd := &cobra.Command{
		Use:          "valor",
		Short:        "valor [flags]",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			rcp, err := loadRecipe(path, defaultRecipeType, defaultRecipeFormat)
			if err != nil {
				return errors.New(string(err.JSON()))
			}
			pipeline, err := core.NewPipeline(rcp, batch, progressType)
			if err != nil {
				return errors.New(string(err.JSON()))
			}
			if err := pipeline.Execute(); err != nil {
				return errors.New(string(err.JSON()))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "recipe", "r", defaultRecipePath, "path of the recipe")
	cmd.Flags().IntVarP(&batch, "batch", "b", defaultBatchSize, "batch size for each process")
	cmd.Flags().StringVarP(&progressType, "progress", "p", defaultProgressType, "progress of the process to show")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadRecipe(path, _type, format string) (*recipe.Recipe, model.Error) {
	fnReader, err := io.Readers.Get(_type)
	if err != nil {
		return nil, err
	}
	getPath := func() string {
		return path
	}
	filterPath := func(p string) bool {
		return true
	}
	postProcess := func(p string, c []byte) (*model.Data, model.Error) {
		return &model.Data{
			Content: c,
			Path:    p,
			Type:    format,
		}, nil
	}
	reader := fnReader(getPath, filterPath, postProcess)
	decode, err := endec.Decodes.Get(format)
	if err != nil {
		return nil, err
	}
	rcp, err := recipe.Load(reader, decode)
	if err != nil {
		return nil, err
	}
	if err := validator.New().Struct(rcp); err != nil {
		return nil, model.BuildError(path, err)
	}
	return rcp, nil
}
