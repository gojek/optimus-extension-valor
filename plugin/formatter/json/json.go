package json

import (
	"bytes"
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/registry/formatter"

	"gopkg.in/yaml.v3"
)

// ToJSON formats input from JSON to JSON, which does nothing
func ToJSON(input []byte) ([]byte, error) {
	return input, nil
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
	formatter.Formatters.Register("json", "json", ToJSON)
	formatter.Formatters.Register("json", "yaml", ToYAML)
}
