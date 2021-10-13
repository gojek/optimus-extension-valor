package json

import (
	"encoding/json"

	"github.com/gojek/optimus-extension-valor/plugin"
	"github.com/gojek/optimus-extension-valor/registry/endec"
)

const _type = "json"

// NewEncoder initializes JSON encoding function
func NewEncoder() plugin.Encoder {
	return json.Marshal
}

// NewDecoder initializes JSON decodign function
func NewDecoder() plugin.Decoder {
	return json.Unmarshal
}

func init() {
	endec.Encoders.Register(_type, NewEncoder)
	endec.Decoders.Register(_type, NewDecoder)
}
