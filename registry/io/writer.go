package io

import (
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/plugin"
)

// WriteFn is a getter for IO Writer instance
type WriteFn func(path string, metadata map[string]string) plugin.Writer

// WriterFactory is a factory for Writer
type WriterFactory struct {
	typeToFn map[string]WriteFn
}

// Register registers a factory function for a type
func (w *WriterFactory) Register(_type string, fn WriteFn) error {
	_type = strings.ToLower(_type)
	if w.typeToFn[_type] != nil {
		return fmt.Errorf("[%s] is already registered", _type)
	}
	w.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a type
func (w *WriterFactory) Get(_type string) (WriteFn, error) {
	_type = strings.ToLower(_type)
	if w.typeToFn[_type] == nil {
		return nil, fmt.Errorf("[%s] is not registered", _type)
	}
	return w.typeToFn[_type], nil
}

// NewWriterFactory initializes factory Writer
func NewWriterFactory() *WriterFactory {
	return &WriterFactory{
		typeToFn: make(map[string]WriteFn),
	}
}
