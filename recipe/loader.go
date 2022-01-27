package recipe

import (
	"bytes"
	"errors"

	"github.com/gojek/optimus-extension-valor/model"
)

// Load loads recipe from the passed Reader with Decoder to decode
func Load(reader model.Reader, decode model.Decode) (*Recipe, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}
	if decode == nil {
		return nil, errors.New("decode is nil")
	}
	data, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data.Content)) == 0 {
		return nil, errors.New("content is empty")
	}
	output := &Recipe{}
	err = decode(data.Content, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
