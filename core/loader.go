package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

// DefinitionWrapper is a wrapper for definition
type DefinitionWrapper struct {
	Definition *model.Definition
	Error      model.Error
}

// SchemaWrapper is a wrapper for schema
type SchemaWrapper struct {
	Schema *model.Schema
	Error  model.Error
}

// ProcedureWrapper is a wrapper for procedure
type ProcedureWrapper struct {
	Procedure *model.Procedure
	Error     model.Error
}

// Loader is a loader management
type Loader struct {
}

// LoadResource loads a Resource based on its recipe
func (l *Loader) LoadResource(rcp *recipe.Resource) (*model.Resource, model.Error) {
	const defaultErrKey = "LoadResource"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("resource recipe is nil"))
	}
	listOfData, err := l.loadAllData(rcp.Path, rcp.Type, rcp.Format)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return &model.Resource{
		Name:           rcp.Name,
		ListOfData:     listOfData,
		FrameworkNames: rcp.FrameworkNames,
	}, nil
}

// LoadFramework loads a framework based on its recipe
func (l *Loader) LoadFramework(rcp *recipe.Framework) (*model.Framework, model.Error) {
	const defaultErrKey = "LoadFramework"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("framework recipe is nil"))
	}
	definitionChans := l.DispatchDefinitionChans(rcp.Definitions)
	schemaChans := l.DispatchSchemaChans(rcp.Schemas)
	procedureChans := l.DispatchProcedureChans(rcp.Procedures)

	definitions, defError := l.acceptDefinitionChans(definitionChans)
	schemas, schError := l.acceptSchemaChans(schemaChans)
	procedures, proError := l.acceptProcedureChans(procedureChans)

	errOutput := model.CombineErrors(defError, schError, proError)
	if len(errOutput) > 0 {
		return nil, errOutput
	}
	return &model.Framework{
		Name:          rcp.Name,
		Definitions:   definitions,
		Schemas:       schemas,
		Procedures:    procedures,
		OutputTargets: l.convertOutputTargets(rcp.OutputTargets),
	}, nil
}

func (l *Loader) convertOutputTargets(targets []*recipe.OutputTarget) []*model.OutputTarget {
	output := make([]*model.OutputTarget, len(targets))
	for i, t := range targets {
		output[i] = &model.OutputTarget{
			Name:   t.Name,
			Format: t.Format,
			Type:   t.Type,
			Path:   t.Path,
		}
	}
	return output
}

func (l *Loader) acceptDefinitionChans(definitionChans []chan *DefinitionWrapper) ([]*model.Definition, model.Error) {
	const defaultErrKey = "acceptDefinitionChans"
	definitions := make([]*model.Definition, len(definitionChans))
	errOutput := make(model.Error)
	for i, ch := range definitionChans {
		defWrapper := <-ch
		if defWrapper.Error != nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			errOutput[key] = defWrapper.Error
			continue
		}
		definitions[i] = defWrapper.Definition
	}
	return definitions, errOutput
}

func (l *Loader) acceptSchemaChans(schemaChans []chan *SchemaWrapper) ([]*model.Schema, model.Error) {
	const defaultErrKey = "acceptSchemaChans"
	schemas := make([]*model.Schema, len(schemaChans))
	errOutput := make(model.Error)
	for i, ch := range schemaChans {
		schWrapper := <-ch
		if schWrapper.Error != nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			errOutput[key] = schWrapper.Error
			continue
		}
		schemas[i] = schWrapper.Schema
	}
	return schemas, errOutput
}

func (l *Loader) acceptProcedureChans(procedureChans []chan *ProcedureWrapper) ([]*model.Procedure, model.Error) {
	const defaultErrKey = "acceptProcedureChans"
	procedures := make([]*model.Procedure, len(procedureChans))
	errOutput := make(model.Error)
	for i, ch := range procedureChans {
		proWrapper := <-ch
		if proWrapper.Error != nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			errOutput[key] = proWrapper.Error
			continue
		}
		procedures[i] = proWrapper.Procedure
	}
	if len(errOutput) > 0 {
		return nil, errOutput
	}
	return procedures, nil
}

// DispatchProcedureChans dispatches channel to load procedures
func (l *Loader) DispatchProcedureChans(procedureRcps []*recipe.Procedure) []chan *ProcedureWrapper {
	procedureChans := make([]chan *ProcedureWrapper, len(procedureRcps))
	for i, procedureRcp := range procedureRcps {
		ch := make(chan *ProcedureWrapper)
		go func(c chan *ProcedureWrapper, r *recipe.Procedure) {
			procedure, err := l.LoadProcedure(r)
			if err != nil {
				c <- &ProcedureWrapper{
					Error: err,
				}
			} else {
				c <- &ProcedureWrapper{
					Procedure: procedure,
				}

			}
		}(ch, procedureRcp)
		procedureChans[i] = ch
	}
	return procedureChans
}

// DispatchSchemaChans dispatches channel to load schemas
func (l *Loader) DispatchSchemaChans(schemaRcps []*recipe.Schema) []chan *SchemaWrapper {
	schemaChans := make([]chan *SchemaWrapper, len(schemaRcps))
	for i, schemaRcp := range schemaRcps {
		ch := make(chan *SchemaWrapper)
		go func(c chan *SchemaWrapper, r *recipe.Schema) {
			schema, err := l.LoadSchema(r)
			if err != nil {
				c <- &SchemaWrapper{
					Error: err,
				}
			} else {
				c <- &SchemaWrapper{
					Schema: schema,
				}

			}
		}(ch, schemaRcp)
		schemaChans[i] = ch
	}
	return schemaChans
}

// DispatchDefinitionChans dispatches channel to load definitions
func (l *Loader) DispatchDefinitionChans(definitionRcps []*recipe.Definition) []chan *DefinitionWrapper {
	definitionChans := make([]chan *DefinitionWrapper, len(definitionRcps))
	for i, definitionRcp := range definitionRcps {
		ch := make(chan *DefinitionWrapper)
		go func(c chan *DefinitionWrapper, r *recipe.Definition) {
			definition, err := l.LoadDefinition(r)
			if err != nil {
				c <- &DefinitionWrapper{
					Error: err,
				}
			} else {
				c <- &DefinitionWrapper{
					Definition: definition,
				}

			}
		}(ch, definitionRcp)
		definitionChans[i] = ch
	}
	return definitionChans
}

// LoadDefinition loads definition based on its recipe
func (l *Loader) LoadDefinition(rcp *recipe.Definition) (*model.Definition, model.Error) {
	const defaultErrKey = "LoadDefinition"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("definition recipe is nil"))
	}
	listOfData, err := l.loadAllData(rcp.Path, rcp.Type, rcp.Format)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	var functionData *model.Data
	if rcp.Function != nil {
		data, err := l.loadOneData(rcp.Function.Path, rcp.Function.Type, jsonnetFormat)
		if err != nil {
			return nil, model.BuildError(defaultErrKey, err)
		}
		functionData = data
	}
	return &model.Definition{
		Name:         rcp.Name,
		ListOfData:   listOfData,
		FunctionData: functionData,
	}, nil
}

// LoadSchema loads schema based on its recipe
func (l *Loader) LoadSchema(rcp *recipe.Schema) (*model.Schema, model.Error) {
	const defaultErrKey = "LoadSchema"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("schema recipe is nil"))
	}
	data, err := l.loadOneData(rcp.Path, rcp.Type, jsonFormat)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return &model.Schema{
		Name: rcp.Name,
		Data: data,
	}, nil
}

// LoadProcedure loads procedure based on its recipe
func (l *Loader) LoadProcedure(rcp *recipe.Procedure) (*model.Procedure, model.Error) {
	const defaultErrKey = "LoadProcedure"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("procedure recipe is nil"))
	}
	data, err := l.loadOneData(rcp.Path, rcp.Type, jsonnetFormat)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return &model.Procedure{
		Name:          rcp.Name,
		OutputIsError: rcp.OutputIsError,
		Data:          data,
	}, nil
}

func (l *Loader) loadOneData(path, _type, format string) (*model.Data, model.Error) {
	reader, err := l.getLoadReader(path, _type, format)
	if err != nil {
		return nil, err
	}
	return reader.ReadOne()
}

func (l *Loader) loadAllData(path, _type, format string) ([]*model.Data, model.Error) {
	const defaultErrKey = "loadAllData"
	reader, err := l.getLoadReader(path, _type, format)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	data, err := reader.ReadAll()
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return data, nil
}

func (l *Loader) getLoadReader(path, _type, format string) (model.Reader, model.Error) {
	const defaultErrKey = "getLoadReader"
	readerFn, err := io.Readers.Get(_type)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	reader := readerFn(
		l.getPath(path),
		l.filterPath(format),
		l.postProcess(_type, format),
	)
	return reader, nil
}

func (l *Loader) getPath(path string) model.GetPath {
	return func() string {
		return path
	}
}

func (l *Loader) filterPath(suffix string) model.FilterPath {
	return func(path string) bool {
		return strings.HasSuffix(path, suffix)
	}
}

func (l *Loader) postProcess(_type, format string) model.PostProcess {
	return func(path string, content []byte) (*model.Data, model.Error) {
		if !skipReformat[format] {
			fn, err := formatter.Formats.Get(format, jsonFormat)
			if err != nil {
				return nil, err
			}
			reformattedContent, err := fn(content)
			if err != nil {
				return nil, err
			}
			content = reformattedContent
		}
		return &model.Data{
			Path:    path,
			Type:    _type,
			Content: content,
		}, nil
	}
}
