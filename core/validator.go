package core

import (
	"errors"
	"fmt"
	"sync"

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
	outputError := make(model.Error)
	for i, sch := range framework.Schemas {
		if sch == nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			outputError[key] = errors.New("schema is nil")
		}
	}
	if len(outputError) > 0 {
		return nil, outputError
	}
	return &Validator{
		framework: framework,
	}, nil
}

// Validate validates a Resource data against all schemas
func (v *Validator) Validate(resourceData *model.Data) model.Error {
	const defaultErrKey = "Validate"
	if resourceData == nil {
		return model.BuildError(defaultErrKey, errors.New("resource data is nil"))
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputError := make(model.Error)
	for i, schema := range v.framework.Schemas {
		wg.Add(1)
		func(idx int, w *sync.WaitGroup, m *sync.Mutex, sch *model.Schema, rsc *model.Data) {
			defer w.Done()

			if err := v.validateResourceToSchema(sch.Data, rsc); err != nil {
				key := fmt.Sprintf("%s [%s: %d]", defaultErrKey, sch.Name, idx)
				m.Lock()
				outputError[key] = err
				m.Unlock()
			}
		}(i, wg, mtx, schema, resourceData)
	}
	if len(outputError) > 0 {
		return outputError
	}
	return nil
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
	outputError := make(model.Error)
	for _, r := range result.Errors() {
		field := r.Field()
		msg := r.Description()
		outputError[field] = msg
	}
	if len(outputError) > 0 {
		return outputError
	}
	return nil
}
