package core

import (
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

// LoadResource loads the content of a Resource based on its recipe
func LoadResource(rcp *recipe.Resource) (*Resource, error) {
	container, err := LoadContainer(&recipe.Metadata{
		Name:   rcp.Name,
		Format: rcp.Format,
		Type:   rcp.Type,
		Path:   rcp.Path,
	})
	if err != nil {
		return nil, err
	}
	return &Resource{
		container:      container,
		frameworkNames: rcp.FrameworkNames,
	}, nil
}

// LoadFramework loads a framework based on its recipe
func LoadFramework(rcp *recipe.Framework) (*Framework, error) {
	schemas, err := loadContainerList(rcp.Schemas)
	if err != nil {
		return nil, err
	}
	definitions, err := loadContainerList(rcp.Definitions)
	if err != nil {
		return nil, err
	}
	procedures, err := loadContainerList(rcp.Procedures)
	if err != nil {
		return nil, err
	}
	return &Framework{
		allowError:    rcp.AllowError,
		name:          rcp.Name,
		schemas:       schemas,
		definitions:   definitions,
		procedures:    procedures,
		outputTargets: rcp.OutputTargets,
	}, nil
}

func loadContainerList(rcp []*recipe.Metadata) ([]*Container, error) {
	var containers []*Container
	for _, c := range rcp {
		ctainer, err := LoadContainer(c)
		if err != nil {
			return nil, err
		}
		containers = append(containers, ctainer)
	}
	return containers, nil
}

// LoadContainer loads a container based on its type, path, and format
func LoadContainer(rcp *recipe.Metadata) (*Container, error) {
	fn, err := io.Readers.Get(rcp.Type)
	if err != nil {
		return nil, err
	}
	newFormat := rcp.Format
	if !isNotToFormat[newFormat] {
		newFormat = defaultFormat
	}
	reader := fn(rcp.Path, map[string]string{
		"name":   rcp.Name,
		"format": newFormat,
		"type":   rcp.Type,
		"path":   rcp.Path,
	})
	var dataList []*model.Data
	for reader.Next() {
		data, err := reader.Read()
		if err != nil {
			return nil, err
		}
		if !isNotToFormat[rcp.Format] {
			content, err := FormatContent(rcp.Format, defaultFormat, data.Content)
			if err != nil {
				return nil, err
			}
			data.Content = content
		}
		dataList = append(dataList, data)
	}
	return &Container{
		name:          rcp.Name,
		data:          dataList,
		outputIsError: rcp.OutputIsError,
	}, nil
}
