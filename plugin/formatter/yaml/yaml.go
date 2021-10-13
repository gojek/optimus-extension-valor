package yaml

import (
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/registry/formatter"

	"gopkg.in/yaml.v3"
)

// ToJSON formats input from YAML to JSON
func ToJSON(input []byte) ([]byte, error) {
	var t interface{}
	err := yaml.Unmarshal(input, &t)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(t, "", "  ")
}

func init() {
	formatter.Formatters.Register("yaml", "json", ToJSON)
}
