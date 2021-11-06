package endec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
)

// DecodeFactory is a factory for Decode
type DecodeFactory struct {
	typeToFn map[string]model.Decode
}

// Register registers a factory function for a type
func (d *DecodeFactory) Register(format string, fn model.Decode) model.Error {
	const defaultErrKey = "Register"
	if fn == nil {
		return model.BuildError(defaultErrKey, errors.New("Decode is nil"))
	}
	format = strings.ToLower(format)
	if d.typeToFn[format] != nil {
		return model.BuildError(defaultErrKey, fmt.Errorf("[%s] is already registered", format))
	}
	d.typeToFn[format] = fn
	return nil
}

// Get gets a factory function based on a type
func (d *DecodeFactory) Get(format string) (model.Decode, model.Error) {
	const defaultErrKey = "Get"
	format = strings.ToLower(format)
	if d.typeToFn[format] == nil {
		return nil, model.BuildError(defaultErrKey, fmt.Errorf("[%s] is not registered", format))
	}
	return d.typeToFn[format], nil
}

// NewDecodeFactory initializes factory Decode
func NewDecodeFactory() *DecodeFactory {
	return &DecodeFactory{
		typeToFn: make(map[string]model.Decode),
	}
}
