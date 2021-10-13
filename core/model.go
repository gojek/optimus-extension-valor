package core

import (
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
)

// Container contains one or more data
type Container struct {
	name          string
	data          []*model.Data
	outputIsError bool
}

// Framework represents a Framework content
type Framework struct {
	allowError    bool
	name          string
	schemas       []*Container
	definitions   []*Container
	procedures    []*Container
	outputTargets []*recipe.Metadata
}

// Resource represents a Resource content
type Resource struct {
	container      *Container
	frameworkNames []string
}

type schemaResult struct {
	record *model.Data
	schema *model.Data
	result []byte
	err    error

	outputTargets []*recipe.Metadata
}

type procedureResult struct {
	record  *model.Data
	snippet string
	result  []byte
	err     error

	outputTargets []*recipe.Metadata
}
