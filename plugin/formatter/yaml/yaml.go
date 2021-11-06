package yaml

import (
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/formatter"

	"gopkg.in/yaml.v3"
)

// ToJSON formats input from YAML to JSON
func ToJSON(input []byte) ([]byte, model.Error) {
	const defaultErrKey = "ToJSON"
	var t interface{}
	err := yaml.Unmarshal(input, &t)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	output, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return output, nil
}

func init() {
	err := formatter.Formats.Register("yaml", "json", ToJSON)
	if err != nil {
		panic(err)
	}
}
