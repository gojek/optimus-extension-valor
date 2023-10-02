package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

// Loader is a loader management
type Loader struct {
}

// LoadFramework loads framework based on the specified recipe
func (l *Loader) LoadFramework(rcp *recipe.Framework) (*model.Framework, error) {
	if rcp == nil {
		return nil, errors.New("framework recipe is nil")
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
		Name:        rcp.Name,
		Definitions: definitions,
		Schemas:     schemas,
		Procedures:  procedures,
	}, nil
}

func (l *Loader) loadAllProcedures(rcps []*recipe.Procedure) ([]*model.Procedure, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Procedure, len(rcps))
	outputError := &model.Error{}
	for i, rcp := range rcps {
		wg.Add(1)
		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, r *recipe.Procedure) {
			defer wg.Done()

			procedure, err := l.LoadProcedure(r)
			if err != nil {
				key := fmt.Sprintf("%d", idx)
				if r != nil {
					key = r.Name
				}
				outputError.Add(key, err)
			} else {
				m.Lock()
				outputData[idx] = procedure
				m.Unlock()
			}
		}(i, wg, mtx, rcp)
	}
	wg.Wait()

	if outputError.Length() > 0 {
		return nil, outputError
	}
	return outputData, nil
}

// LoadProcedure loads procedure based on the specified recipe
func (l *Loader) LoadProcedure(rcp *recipe.Procedure) (*model.Procedure, error) {
	if rcp == nil {
		return nil, errors.New("procedure recipe is nil")
	}
	paths, err := ExplorePaths(rcp.Path, rcp.Type, jsonnetFormat)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("[%s] procedure for recipe [%s] cannot be found", jsonnetFormat, rcp.Name)
	}
	data, err := l.LoadData(paths[0], rcp.Type, jsonnetFormat)
	if err != nil {
		return nil, err
	}
	return &model.Procedure{
		Name:   rcp.Name,
		Data:   data,
		Output: l.convertOutput(rcp.Output),
	}, nil
}

func (l *Loader) loadAllSchemas(rcps []*recipe.Schema) ([]*model.Schema, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Schema, len(rcps))
	outputError := &model.Error{}
	for i, rcp := range rcps {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, r *recipe.Schema) {
			defer wg.Done()
			schema, err := l.LoadSchema(r)
			if err != nil {
				key := fmt.Sprintf("%d", idx)
				if r != nil {
					key = r.Name
				}
				outputError.Add(key, err)
			} else {
				m.Lock()
				outputData[idx] = schema
				m.Unlock()
			}
		}(i, wg, mtx, rcp)
	}
	wg.Wait()

	if outputError.Length() > 0 {
		return nil, outputError
	}
	return outputData, nil
}

// LoadSchema loads schema based on the specified recipe
func (l *Loader) LoadSchema(rcp *recipe.Schema) (*model.Schema, error) {
	if rcp == nil {
		return nil, errors.New("schema recipe is nil")
	}
	paths, err := ExplorePaths(rcp.Path, rcp.Type, jsonFormat)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("[%s] schema for recipe [%s] cannot be found", jsonFormat, rcp.Name)
	}
	data, err := l.LoadData(paths[0], rcp.Type, jsonFormat)
	if err != nil {
		return nil, err
	}
	return &model.Schema{
		Name:   rcp.Name,
		Data:   data,
		Output: l.convertOutput(rcp.Output),
	}, nil
}

func (l *Loader) convertOutput(output *recipe.Output) *model.Output {
	if output == nil {
		return nil
	}
	targets := make([]*model.Target, len(output.Targets))
	for i, t := range output.Targets {
		targets[i] = &model.Target{
			Name:   t.Name,
			Format: t.Format,
			Type:   t.Type,
			Path:   t.Path,
		}
	}
	return &model.Output{
		TreatAs: model.OutputTreatment(output.TreatAs),
		Targets: targets,
	}
}

func (l *Loader) loadAllDefinitions(rcps []*recipe.Definition) ([]*model.Definition, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputData := make([]*model.Definition, len(rcps))
	outputError := &model.Error{}
	for i, rcp := range rcps {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex, r *recipe.Definition) {
			defer wg.Done()
			definition, err := l.LoadDefinition(r)
			if err != nil {
				key := fmt.Sprintf("%d", idx)
				if r != nil {
					key = r.Name
				}
				outputError.Add(key, err)
			} else {
				m.Lock()
				outputData[idx] = definition
				m.Unlock()
			}
		}(i, wg, mtx, rcp)
	}
	wg.Wait()

	if outputError.Length() > 0 {
		return nil, outputError
	}
	return outputData, nil
}

// LoadDefinition loads definition based on the specified recipe
func (l *Loader) LoadDefinition(rcp *recipe.Definition) (*model.Definition, error) {
	if rcp == nil {
		return nil, errors.New("definition recipe is nil")
	}
	paths, err := ExplorePaths(rcp.Path, rcp.Type, rcp.Format)
	if err != nil {
		return nil, err
	}
	listOfData := make([]*model.Data, len(paths))
	for i, p := range paths {
		data, err := l.LoadData(p, rcp.Type, rcp.Format)
		if err != nil {
			return nil, err
		}
		listOfData[i] = data
	}
	var functionData *model.Data
	if rcp.Function != nil {
		data, err := l.LoadData(rcp.Function.Path, rcp.Function.Type, jsonnetFormat)
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

// LoadData loads data based on the specified path, type, and format
func (l *Loader) LoadData(path, _type, format string) (*model.Data, error) {
	reader, err := l.getReader(path, _type, format)
	if err != nil {
		return nil, err
	}
	return reader.Read()
}

func (l *Loader) getReader(path, _type, format string) (model.Reader, error) {
	readerFn, err := io.Readers.Get(_type)
	if err != nil {
		return nil, err
	}
	reader := readerFn(
		func() string {
			return path
		},
		func(path string, content []byte) (*model.Data, error) {
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
		},
	)
	return reader, nil
}
