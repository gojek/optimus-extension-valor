package yaml

import (
	"github.com/gojek/optimus-extension-valor/plugin"
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"gopkg.in/yaml.v3"
)

const _type = "yaml"

// NewEncoder initializes YAML encoder function
func NewEncoder() plugin.Encoder {
	return yaml.Marshal
}

// NewDecoder initializes YAML decoder function
func NewDecoder() plugin.Decoder {
	return yaml.Unmarshal
}

func init() {
	endec.Decoders.Register(_type, NewDecoder)
	endec.Encoders.Register(_type, NewEncoder)
}
