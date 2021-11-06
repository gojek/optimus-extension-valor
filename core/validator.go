package core

import (
	"errors"

	"github.com/gojek/optimus-extension-valor/model"

	"github.com/xeipuuv/gojsonschema"
)

// Validator is a validator for Resource against a Schema
type Validator struct {
	framework *model.Framework
}

// NewValidator initializes Validator
func NewValidator(framework *model.Framework) (*Validator, model.Error) {
	const defaultErrKey = "NewValidator"
	if framework == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("framework is nil"))
	}
	return &Validator{
		framework: framework,
	}, nil
}

// Validate validates a Resource data against all schemas
func (v *Validator) Validate(resourceData *model.Data) model.Error {
	const defaultErrKey = "Validate"
	validateChans := v.dispatchValidate(v.framework.Schemas, resourceData)
	results := make([]model.Error, len(validateChans))
	for i, ch := range validateChans {
		results[i] = <-ch
	}
	return model.CombineErrors(results...)
}

func (v *Validator) dispatchValidate(schemaList []*model.Schema, resourceData *model.Data) []chan model.Error {
	const defaultErrKey = "dispatchValidate"
	validateChans := make([]chan model.Error, len(schemaList))
	for i, schema := range schemaList {
		ch := make(chan model.Error)
		go func(c chan model.Error, sc *model.Schema, rsc *model.Data) {
			if sc == nil {
				c <- model.BuildError(defaultErrKey, errors.New("schema is nil"))
				return
			}
			c <- v.validateResourceToSchema(sc.Data, rsc)
		}(ch, schema, resourceData)
		validateChans[i] = ch
	}
	return validateChans
}

func (v *Validator) validateResourceToSchema(schemaData *model.Data, resourceData *model.Data) model.Error {
	const defaultErrKey = "validateResourceToSchema"
	if schemaData == nil {
		return model.BuildError(defaultErrKey, errors.New("schema data is nil"))
	}
	if resourceData == nil {
		return model.BuildError(defaultErrKey, errors.New("resource data is nil"))
	}
	schemaLoader := gojsonschema.NewStringLoader(string(schemaData.Content))
	recordLoader := gojsonschema.NewStringLoader(string(resourceData.Content))
	result, validateErr := gojsonschema.Validate(schemaLoader, recordLoader)
	if validateErr != nil {
		return model.BuildError(defaultErrKey, validateErr)
	}
	if result.Valid() {
		return nil
	}
	output := make([]model.Error, len(result.Errors()))
	for i, r := range result.Errors() {
		field := r.Field()
		msg := r.Description()
		output[i] = model.BuildError(field, msg)
	}
	return model.CombineErrors(output...)
}
