package endec

import (
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/plugin"
)

// EncoderFn is a getter for Encoder instance
type EncoderFn func() plugin.Encoder

// EncoderFactory is a factory for Encoder
type EncoderFactory struct {
	typeToFn map[string]EncoderFn
}

// Register registers a factory function for a type
func (e *EncoderFactory) Register(_type string, fn EncoderFn) error {
	_type = strings.ToLower(_type)
	if e.typeToFn[_type] != nil {
		return fmt.Errorf("[%s] is already registered", _type)
	}
	e.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a type
func (e *EncoderFactory) Get(_type string) (EncoderFn, error) {
	_type = strings.ToLower(_type)
	if e.typeToFn[_type] == nil {
		return nil, fmt.Errorf("[%s] is not registered", _type)
	}
	return e.typeToFn[_type], nil
}

// NewEncoderFactory initializes factory Encoder
func NewEncoderFactory() *EncoderFactory {
	return &EncoderFactory{
		typeToFn: make(map[string]EncoderFn),
	}
}
