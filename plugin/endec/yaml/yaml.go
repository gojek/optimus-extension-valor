package yaml

import (
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"gopkg.in/yaml.v3"
)

const format = "yaml"

// NewEncode initializes YAML encoder function
func NewEncode() model.Encode {
	const defaultErrKey = "NewEncode"
	return func(i interface{}) ([]byte, model.Error) {
		output, err := yaml.Marshal(i)
		if err != nil {
			return nil, model.BuildError(defaultErrKey, err)
		}
		return output, nil
	}
}

// NewDecode initializes YAML decoder function
func NewDecode() model.Decode {
	const defaultErrKey = "NewDecode"
	return func(b []byte, i interface{}) model.Error {
		if err := yaml.Unmarshal(b, i); err != nil {
			return model.BuildError(defaultErrKey, err)
		}
		return nil
	}
}

func init() {
	err := endec.Decodes.Register(format, NewDecode())
	if err != nil {
		panic(err)
	}
	err = endec.Encodes.Register(format, NewEncode())
	if err != nil {
		panic(err)
	}
}
