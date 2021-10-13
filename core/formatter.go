package core

import "github.com/gojek/optimus-extension-valor/registry/formatter"

// FormatContent formats a content from a specified format to the target format
func FormatContent(source, target string, content []byte) ([]byte, error) {
	fn, err := formatter.Formatters.Get(source, target)
	if err != nil {
		return nil, err
	}
	return fn(content)
}
