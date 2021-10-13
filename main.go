package main

import (
	"github.com/gojek/optimus-extension-valor/core"
	_ "github.com/gojek/optimus-extension-valor/plugin/endec"
	_ "github.com/gojek/optimus-extension-valor/plugin/formatter"
	_ "github.com/gojek/optimus-extension-valor/plugin/io"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/endec"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/go-playground/validator/v10"
)

var defaultPath = "./valor.yaml"

func main() {
	rcp := loadRecipe()
	eval, err := core.NewPipeline(rcp)
	if err != nil {
		panic(err)
	}
	err = eval.Load()
	if err != nil {
		panic(err)
	}
	err = eval.Build()
	if err != nil {
		panic(err)
	}
	err = eval.Execute()
	if err != nil {
		panic(err)
	}
	err = eval.Flush()
	if err != nil {
		panic(err)
	}
}

func loadRecipe() *recipe.Recipe {
	readerType := "file"
	fnReader, err := io.Readers.Get(readerType)
	if err != nil {
		panic(err)
	}
	reader := fnReader(defaultPath, nil)
	decoderType := "yaml"
	fnDecoder, err := endec.Decoders.Get(decoderType)
	if err != nil {
		panic(err)
	}
	decoder := fnDecoder()
	rcp, err := recipe.LoadWithReader(reader, decoder)
	if err != nil {
		panic(err)
	}
	if err := validator.New().Struct(rcp); err != nil {
		panic(err)
	}
	return rcp
}
