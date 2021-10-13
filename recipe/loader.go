package recipe

import (
	"bytes"
	"errors"

	"github.com/gojek/optimus-extension-valor/plugin"
)

// LoadWithReader loads recipe from the passed Reader with Decoder to decode
func LoadWithReader(reader plugin.Reader, decoder plugin.Decoder) (*Recipe, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}
	if decoder == nil {
		return nil, errors.New("decoder is nil")
	}
	data, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data.Content)) == 0 {
		return nil, errors.New("content is empty")
	}
	output := &Recipe{}
	err = decoder(data.Content, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
