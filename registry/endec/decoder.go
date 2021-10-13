package endec

import (
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/plugin"
)

// DecoderFn is a getter for Decoder instance
type DecoderFn func() plugin.Decoder

// DecoderFactory is a factory for Decoder
type DecoderFactory struct {
	typeToFn map[string]DecoderFn
}

// Register registers a factory function for a type
func (d *DecoderFactory) Register(_type string, fn DecoderFn) error {
	_type = strings.ToLower(_type)
	if d.typeToFn[_type] != nil {
		return fmt.Errorf("[%s] is already registered", _type)
	}
	d.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a type
func (d *DecoderFactory) Get(_type string) (DecoderFn, error) {
	_type = strings.ToLower(_type)
	if d.typeToFn[_type] == nil {
		return nil, fmt.Errorf("[%s] is not registered", _type)
	}
	return d.typeToFn[_type], nil
}

// NewDecoderFactory initializes factory Decoder
func NewDecoderFactory() *DecoderFactory {
	return &DecoderFactory{
		typeToFn: make(map[string]DecoderFn),
	}
}
