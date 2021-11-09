package progress

import (
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
)

// Progresses is a factory for Progress
var Progresses = NewFactory()

// Factory is a factory for Progress
type Factory struct {
	typeToFn map[string]model.NewProgress
}

// Register registers a factory function for a specified type
func (f *Factory) Register(_type string, fn model.NewProgress) model.Error {
	const defaultErrKey = "Register"
	_type = strings.ToLower(_type)
	if f.typeToFn[_type] != nil {
		return model.BuildError(defaultErrKey, fmt.Errorf("[%s] is already registered", _type))
	}
	f.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a specified type
func (f *Factory) Get(_type string) (model.NewProgress, model.Error) {
	const defaultErrKey = "Get"
	_type = strings.ToLower(_type)
	if f.typeToFn[_type] == nil {
		return nil, model.BuildError(defaultErrKey, fmt.Errorf("[%s] is not registered", _type))
	}
	return f.typeToFn[_type], nil
}

// NewFactory initializes factory Formatter
func NewFactory() *Factory {
	return &Factory{
		typeToFn: make(map[string]model.NewProgress),
	}
}
