package yaml

import (
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/registry/formatter"

	"gopkg.in/yaml.v3"
)

const (
	jsonType = "json"
	yamlType = "yaml"
)

// ToJSON formats input from YAML to JSON
func ToJSON(input []byte) ([]byte, error) {
	var t interface{}
	err := yaml.Unmarshal(input, &t)
	if err != nil {
		return nil, err
	}
	output, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}

func init() {
	err := formatter.Formats.Register(yamlType, jsonType, ToJSON)
	if err != nil {
		panic(err)
	}
}
