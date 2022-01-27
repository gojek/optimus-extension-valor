package explorer

import (
	"errors"
	"fmt"

	"github.com/gojek/optimus-extension-valor/model"
)

// Explorers is a factory for explorer
var Explorers = NewFactory()

// Factory is a factory for Explorer
type Factory struct {
	typeToFn map[string]model.ExplorePath
}

// Register registers an explorer based on the type
func (f *Factory) Register(_type string, fn model.ExplorePath) error {
	if fn == nil {
		return errors.New("Explorer is nil")
	}
	if f.typeToFn[_type] != nil {
		return fmt.Errorf("[%s] is already registered", _type)
	}
	f.typeToFn[_type] = fn
	return nil
}

// Get gets an explorer based on a specified type
func (f *Factory) Get(_type string) (model.ExplorePath, error) {
	if f.typeToFn[_type] == nil {
		return nil, fmt.Errorf("[%s] is not registered", _type)
	}
	return f.typeToFn[_type], nil
}

// NewFactory initializes factory Formatter
func NewFactory() *Factory {
	return &Factory{
		typeToFn: make(map[string]model.ExplorePath),
	}
}
