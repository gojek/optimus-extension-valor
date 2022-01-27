package io

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
)

// ReaderFn is a getter for IO Reader instance
type ReaderFn func(
	getPath model.GetPath,
	postProcess model.PostProcess,
) model.Reader

// ReaderFactory is a factory for Reader
type ReaderFactory struct {
	typeToFn map[string]ReaderFn
}

// Register registers a factory function for a type
func (r *ReaderFactory) Register(_type string, fn ReaderFn) error {
	if fn == nil {
		return errors.New("ReaderFn is nil")
	}
	_type = strings.ToLower(_type)
	if r.typeToFn[_type] != nil {
		return fmt.Errorf("[%s] is already registered", _type)
	}
	r.typeToFn[_type] = fn
	return nil
}

// Get gets a factory function based on a type
func (r *ReaderFactory) Get(_type string) (ReaderFn, error) {
	_type = strings.ToLower(_type)
	if r.typeToFn[_type] == nil {
		return nil, fmt.Errorf("[%s] is not registered", _type)
	}
	return r.typeToFn[_type], nil
}

// NewReaderFactory initializes factory Reader
func NewReaderFactory() *ReaderFactory {
	return &ReaderFactory{
		typeToFn: make(map[string]ReaderFn),
	}
}
