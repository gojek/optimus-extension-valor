package json

import (
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/endec"
)

const format = "json"

// NewEncode initializes JSON encoding function
func NewEncode() model.Encode {
	return func(i interface{}) ([]byte, error) {
		output, err := json.Marshal(i)
		if err != nil {
			return nil, err
		}
		return output, nil
	}
}

// NewDecode initializes JSON decodign function
func NewDecode() model.Decode {
	return func(b []byte, i interface{}) error {
		if err := json.Unmarshal(b, i); err != nil {
			return err
		}
		return nil
	}
}

func init() {
	err := endec.Encodes.Register(format, NewEncode())
	if err != nil {
		panic(err)
	}
	err = endec.Decodes.Register(format, NewDecode())
	if err != nil {
		panic(err)
	}
}
