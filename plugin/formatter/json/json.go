package json

import (
	"bytes"
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/registry/formatter"

	"gopkg.in/yaml.v3"
)

const (
	jsonType = "json"
	yamlType = "yaml"
)

// ToJSON formats input from JSON to JSON, which does nothing
func ToJSON(input []byte) ([]byte, error) {
	output := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		output[i] = input[i]
	}
	return output, nil
}

// ToYAML formats input from JSON to YAML
func ToYAML(input []byte) ([]byte, error) {
	var t interface{}
	err := json.Unmarshal(input, &t)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	y := yaml.NewEncoder(&b)
	y.SetIndent(2)
	if err := y.Encode(t); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func init() {
	err := formatter.Formats.Register(jsonType, jsonType, ToJSON)
	if err != nil {
		panic(err)
	}
	err = formatter.Formats.Register(jsonType, yamlType, ToYAML)
	if err != nil {
		panic(err)
	}
}
