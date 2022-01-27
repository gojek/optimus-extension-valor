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
func (d *DecodeFactory) Register(format string, fn model.Decode) error {
	if fn == nil {
		return errors.New("Decode is nil")
	}
	format = strings.ToLower(format)
	if d.typeToFn[format] != nil {
		return fmt.Errorf("[%s] is already registered", format)
	}
	d.typeToFn[format] = fn
	return nil
}

// Get gets a factory function based on a type
func (d *DecodeFactory) Get(format string) (model.Decode, error) {
	format = strings.ToLower(format)
	if d.typeToFn[format] == nil {
		return nil, fmt.Errorf("[%s] is not registered", format)
	}
	return d.typeToFn[format], nil
}

// NewDecodeFactory initializes factory Decode
func NewDecodeFactory() *DecodeFactory {
	return &DecodeFactory{
		typeToFn: make(map[string]model.Decode),
	}
}
