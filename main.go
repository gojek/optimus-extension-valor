package main

import (
	"os"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/model"
	_ "github.com/gojek/optimus-extension-valor/plugin/endec"
	_ "github.com/gojek/optimus-extension-valor/plugin/formatter"
	_ "github.com/gojek/optimus-extension-valor/plugin/io"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/endec"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/go-playground/validator/v10"
)

const (
	defaultPath = "./valor.yaml"
	decoderType = "yaml"
)

func main() {
	var args []string
	if len(os.Args) > 1 {
		args = append(args, os.Args[1])
	}
	rcp := loadRecipe()
	writer := getErrorWriter()

	eval, err := core.NewPipeline(rcp)
	if err != nil {
		writeError(writer, err)
	}
	err = eval.Execute()
	if err != nil {
		writeError(writer, err)
	}
}

func writeError(writer model.Writer, err model.Error) {
	data := &model.Data{
		Content: err.JSON(),
	}
	writer.Write(data)
	os.Exit(1)
}

func getErrorWriter() model.Writer {
	writer, err := io.Writers.Get("std")
	if err != nil {
		panic(err)
	}
	return writer
}

func loadRecipe() *recipe.Recipe {
	readerType := "file"
	fnReader, err := io.Readers.Get(readerType)
	if err != nil {
		panic(err)
	}
	getPath := func() string {
		return defaultPath
	}
	filterPath := func(path string) bool {
		return true
	}
	postProcess := func(path string, content []byte) (*model.Data, model.Error) {
		return &model.Data{
			Content: content,
			Path:    path,
			Type:    readerType,
		}, nil
	}
	reader := fnReader(getPath, filterPath, postProcess)
	decode, err := endec.Decodes.Get(decoderType)
	if err != nil {
		panic(err)
	}
	rcp, err := recipe.Load(reader, decode)
	if err != nil {
		panic(err)
	}
	if err := validator.New().Struct(rcp); err != nil {
		panic(err)
	}
	return rcp
}
