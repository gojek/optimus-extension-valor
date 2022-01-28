package core

import (
	"errors"
	"fmt"

	"github.com/gojek/optimus-extension-valor/model"

	"github.com/xeipuuv/gojsonschema"
)

// Validator is a validator for Resource against a Schema
type Validator struct {
	framework *model.Framework
}

// NewValidator initializes Validator
func NewValidator(framework *model.Framework) (*Validator, error) {
	if framework == nil {
		return nil, errors.New("framework is nil")
	}
	outputError := &model.Error{}
	for i, sch := range framework.Schemas {
		if sch == nil {
			key := fmt.Sprintf("%d", i)
			outputError.Add(key, errors.New("schema is nil"))
		}
	}
	if outputError.Length() > 0 {
		return nil, outputError
	}
	return &Validator{
		framework: framework,
	}, nil
}

// Validate validates a Resource data against all schemas
func (v *Validator) Validate(resourceData *model.Data) (bool, error) {
	if resourceData == nil {
		return false, errors.New("resource data is nil")
	}
	for _, schema := range v.framework.Schemas {
		if schema.Data == nil {
			return false, fmt.Errorf("schema data for [%s] is nil", schema.Name)
		}
		schemaLoader := gojsonschema.NewStringLoader(string(schema.Data.Content))
		recordLoader := gojsonschema.NewStringLoader(string(resourceData.Content))
		result, validateErr := gojsonschema.Validate(schemaLoader, recordLoader)
		if validateErr != nil {
			return false, validateErr
		}
		if result.Valid() {
			continue
		}
		businessOutput := &model.Error{}
		for _, r := range result.Errors() {
			field := r.Field()
			msg := r.Description()
			businessOutput.Add(field, msg)
		}
		if businessOutput.Length() > 0 {
			success, err := treatOutput(
				&model.Data{
					Type:    resourceData.Type,
					Path:    resourceData.Path,
					Content: businessOutput.JSON(),
				},
				schema.Output,
			)
			if err != nil {
				return false, err
			}
			if !success {
				return false, nil
			}
		}
	}
	return true, nil
}
