package core

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

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
		return nil, err
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

	definitions, defError := l.loadAllDefinitions(rcp.Definitions)
	if defError != nil {
		return nil, defError
	}
	schemas, schError := l.loadAllSchemas(rcp.Schemas)
	if schError != nil {
		return nil, schError
	}
	procedures, proError := l.loadAllProcedures(rcp.Procedures)
	if proError != nil {
		return nil, proError
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

func (l *Loader) loadAllDefinitions(rcps []*recipe.Definition) ([]*model.Definition, model.Error) {
	const defaultErrKey = "loadAllDefinitions"

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Definition, len(rcps))
	outputError := make(model.Error)
	for i, rcp := range rcps {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, r *recipe.Definition) {
			defer wg.Done()

			definition, err := l.LoadDefinition(r)
			if err != nil {
				key := fmt.Sprintf("%s [%d]", defaultErrKey, idx)
				m.Lock()
				outputError[key] = err
				m.Unlock()
			} else {
				m.Lock()
				outputData[idx] = definition
				m.Unlock()
			}
		}(i, wg, mtx, rcp)
	}
	wg.Wait()

	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputData, nil
}

func (l *Loader) loadAllSchemas(rcps []*recipe.Schema) ([]*model.Schema, model.Error) {
	const defaultErrKey = "loadAllSchemas"

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Schema, len(rcps))
	outputError := make(model.Error)
	for i, rcp := range rcps {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, r *recipe.Schema) {
			defer wg.Done()

			schema, err := l.LoadSchema(r)
			if err != nil {
				key := fmt.Sprintf("%s [%d]", defaultErrKey, idx)
				m.Lock()
				outputError[key] = err
				m.Unlock()
			} else {
				m.Lock()
				outputData[idx] = schema
				m.Unlock()
			}
		}(i, wg, mtx, rcp)
	}
	wg.Wait()

	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputData, nil
}

func (l *Loader) loadAllProcedures(rcps []*recipe.Procedure) ([]*model.Procedure, model.Error) {
	const defaultErrKey = "loadAllProcedures"

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Procedure, len(rcps))
	outputError := make(model.Error)
	for i, rcp := range rcps {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, r *recipe.Procedure) {
			defer wg.Done()

			procedure, err := l.LoadProcedure(r)
			if err != nil {
				key := fmt.Sprintf("%s [%d]", defaultErrKey, idx)
				m.Lock()
				outputError[key] = err
				m.Unlock()
			} else {
				m.Lock()
				outputData[idx] = procedure
				m.Unlock()
			}
		}(i, wg, mtx, rcp)
	}
	wg.Wait()

	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputData, nil
}

// LoadDefinition loads definition based on its recipe
func (l *Loader) LoadDefinition(rcp *recipe.Definition) (*model.Definition, model.Error) {
	const defaultErrKey = "LoadDefinition"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("definition recipe is nil"))
	}
	listOfData, err := l.loadAllData(rcp.Path, rcp.Type, rcp.Format)
	if err != nil {
		return nil, err
	}
	var functionData *model.Data
	if rcp.Function != nil {
		data, err := l.loadOneData(rcp.Function.Path, rcp.Function.Type, jsonnetFormat)
		if err != nil {
			return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *Loader) getLoadReader(path, _type, format string) (model.Reader, model.Error) {
	const defaultErrKey = "getLoadReader"
	readerFn, err := io.Readers.Get(_type)
	if err != nil {
		return nil, err
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
