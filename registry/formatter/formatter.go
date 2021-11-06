package formatter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
)

// Formats is a factory for Format
var Formats = NewFactory()

// Factory is a factory for Formatter
type Factory struct {
	srcToDestToFn map[string]map[string]model.Format
}

// Register registers a factory function for a specified source and destination
func (f *Factory) Register(src, dest string, fn model.Format) model.Error {
	const defaultErrKey = "Register"
	if fn == nil {
		return model.BuildError(defaultErrKey, errors.New("Format is nil"))
	}
	src = strings.ToLower(src)
	dest = strings.ToLower(dest)
	if f.srcToDestToFn[src] != nil && f.srcToDestToFn[src][dest] != nil {
		return model.BuildError(defaultErrKey, fmt.Errorf("[source: %s | target: %s] is already registered", src, dest))
	}
	if f.srcToDestToFn[src] == nil {
		f.srcToDestToFn[src] = make(map[string]model.Format)
	}
	f.srcToDestToFn[src][dest] = fn
	return nil
}

// Get gets a factory function based on a specified source and destination
func (f *Factory) Get(src, dest string) (model.Format, model.Error) {
	const defaultErrKey = "Get"
	src = strings.ToLower(src)
	dest = strings.ToLower(dest)
	if f.srcToDestToFn[src] == nil || f.srcToDestToFn[src][dest] == nil {
		return nil, model.BuildError(defaultErrKey, fmt.Errorf("[source: %s | target: %s] is not registered", src, dest))
	}
	return f.srcToDestToFn[src][dest], nil
}

// NewFactory initializes factory Formatter
func NewFactory() *Factory {
	return &Factory{
		srcToDestToFn: make(map[string]map[string]model.Format),
	}
}
