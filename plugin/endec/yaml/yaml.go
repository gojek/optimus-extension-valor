package yaml

import (
	"fmt"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"gopkg.in/yaml.v3"
)

const format = "yaml"

// NewEncode initializes YAML encoder function
func NewEncode() model.Encode {
	return func(i interface{}) ([]byte, error) {
		var recoverErr error
		output, err := func() ([]byte, error) {
			defer func() {
				if rec := recover(); rec != nil {
					recoverErr = fmt.Errorf("%v", rec)
				}
			}()
			return yaml.Marshal(i)
		}()
		if err != nil {
			return nil, err
		}
		if recoverErr != nil {
			return nil, recoverErr
		}
		return output, nil
	}
}

// NewDecode initializes YAML decoder function
func NewDecode() model.Decode {
	return func(b []byte, i interface{}) error {
		if err := yaml.Unmarshal(b, i); err != nil {
			return err
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
