package io

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
)

// WriterFn is a getter for IO Writer instance
type WriterFn func(treatment model.OutputTreatment) model.Writer

// WriterFactory is a factory for Writer
type WriterFactory struct {
	typeToFn map[string]WriterFn
}

// Register registers a factory function for a type
func (w *WriterFactory) Register(_type string, fn WriterFn) model.Error {
	const defaultErrKey = "Register"
	if fn == nil {
		return model.BuildError(defaultErrKey, errors.New("WriteFn is nil"))
	}
	_type = strings.ToLower(_type)
	if w.typeToFn[_type] != nil {
		return model.BuildError(defaultErrKey, fmt.Errorf("[%s] is already registered", _type))
	}
	w.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a type
func (w *WriterFactory) Get(_type string) (WriterFn, model.Error) {
	const defaultErrKey = "Get"
	_type = strings.ToLower(_type)
	if w.typeToFn[_type] == nil {
		return nil, model.BuildError(defaultErrKey, fmt.Errorf("[%s] is not registered", _type))
	}
	return w.typeToFn[_type], nil
}

// NewWriterFactory initializes factory Writer
func NewWriterFactory() *WriterFactory {
	return &WriterFactory{
		typeToFn: make(map[string]WriterFn),
	}
}
