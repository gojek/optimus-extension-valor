package cmd

import (
	"bytes"
	"errors"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/endec"
	"github.com/gojek/optimus-extension-valor/registry/io"
	"github.com/gojek/optimus-extension-valor/registry/progress"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
)

const defaultProgressType = "progressive"

var progressType string

func getExecuteCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute pipeline based on the specified recipe",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePipeline(recipePath, progressType, enrichWithBatchSize)
		},
	}
	runCmd.PersistentFlags().StringVarP(&recipePath, "recipe-path", "R", defaultRecipePath, "Path of the recipe file")
	runCmd.PersistentFlags().StringVarP(&progressType, "progress-type", "P", defaultProgressType, "Progress type to be used")

	runCmd.AddCommand(getResourceCmd())
	return runCmd
}

func executePipeline(recipePath, progressType string, enrich func(*recipe.Recipe) error) error {
	rcp, err := loadRecipe(recipePath, defaultRecipeType, defaultRecipeFormat)
	if err != nil {
		return err
	}
	if enrich != nil {
		if err := enrich(rcp); err != nil {
			return err
		}
	}
	if err := recipe.Validate(rcp); err != nil {
		return err
	}
	newProgress, err := progress.Progresses.Get(progressType)
	if err != nil {
		return err
	}
	evaluate := getEvaluate()
	pipeline, err := core.NewPipeline(rcp, evaluate, newProgress)
	if err != nil {
		return err
	}
	err = pipeline.Execute()
	if e, ok := err.(*model.Error); ok {
		return errors.New(string(e.JSON()))
	}
	return err
}

func getEvaluate() model.Evaluate {
	vm := jsonnet.MakeVM()
	return func(name, snippet string) (string, error) {
		return vm.EvaluateAnonymousSnippet(name, snippet)
	}
}

func loadRecipe(path, _type, format string) (*recipe.Recipe, error) {
	fnReader, err := io.Readers.Get(_type)
	if err != nil {
		return nil, err
	}
	getPath := func() string {
		return path
	}
	postProcess := func(p string, c []byte) (*model.Data, error) {
		return &model.Data{
			Content: bytes.ToLower(c),
			Path:    p,
			Type:    format,
		}, nil
	}
	reader := fnReader(getPath, postProcess)
	decode, err := endec.Decodes.Get(format)
	if err != nil {
		return nil, err
	}
	return recipe.Load(reader, decode)
}
