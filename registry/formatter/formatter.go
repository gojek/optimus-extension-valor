package formatter

import (
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/plugin"
)

// Formatters is a factory for Formatter
var Formatters = NewFactory()

// Factory is a factory for Formatter
type Factory struct {
	srcToDestToFn map[string]map[string]plugin.Formatter
}

// Register registers a factory function for a specified source and destination
func (r *Factory) Register(src, dest string, fn plugin.Formatter) error {
	src = strings.ToLower(src)
	dest = strings.ToLower(dest)
	if r.srcToDestToFn[src] != nil && r.srcToDestToFn[src][dest] != nil {
		return fmt.Errorf("[source: %s | target: %s] is already registered", src, dest)
	}
	if r.srcToDestToFn[src] == nil {
		r.srcToDestToFn[src] = make(map[string]plugin.Formatter)
	}
	r.srcToDestToFn[src][dest] = fn
	return nil
}

// Get gets a factory function based on a specified source and destination
func (r *Factory) Get(src, dest string) (plugin.Formatter, error) {
	src = strings.ToLower(src)
	dest = strings.ToLower(dest)
	if r.srcToDestToFn[src] == nil || r.srcToDestToFn[src][dest] == nil {
		return nil, fmt.Errorf("[source: %s | target: %s] is not registered", src, dest)
	}
	return r.srcToDestToFn[src][dest], nil
}

// NewFactory initializes factory Formatter
func NewFactory() *Factory {
	return &Factory{
		srcToDestToFn: make(map[string]map[string]plugin.Formatter),
	}
}
