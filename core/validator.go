package core

import (
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateSchema validates record against schema rule
func ValidateSchema(schema, record []byte) ([]byte, error) {
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	recordLoader := gojsonschema.NewBytesLoader(record)
	result, err := gojsonschema.Validate(schemaLoader, recordLoader)
	if err != nil {
		return nil, err
	}
	if result.Valid() {
		return []byte{}, nil
	}
	output := make(map[string][]string)
	for _, r := range result.Errors() {
		field := r.Field()
		msg := r.Description()
		output[field] = append(output[field], msg)
	}
	fn, err := endec.Encoders.Get(defaultFormat)
	if err != nil {
		return nil, err
	}
	encoder := fn()
	return encoder(output)
}
