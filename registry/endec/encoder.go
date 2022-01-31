package endec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
)

// EncodeFactory is a factory for Encode
type EncodeFactory struct {
	typeToFn map[string]model.Encode
}

// Register registers a factory function for a type
func (e *EncodeFactory) Register(_type string, fn model.Encode) error {
	if fn == nil {
		return errors.New("Encode is nil")
	}
	_type = strings.ToLower(_type)
	if e.typeToFn[_type] != nil {
		return fmt.Errorf("[%s] is already registered", _type)
	}
	e.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a type
func (e *EncodeFactory) Get(_type string) (model.Encode, error) {
	_type = strings.ToLower(_type)
	if e.typeToFn[_type] == nil {
		return nil, fmt.Errorf("[%s] is not registered", _type)
	}
	return e.typeToFn[_type], nil
}

// NewEncodeFactory initializes factory Encode
func NewEncodeFactory() *EncodeFactory {
	return &EncodeFactory{
		typeToFn: make(map[string]model.Encode),
	}
}
