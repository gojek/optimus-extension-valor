package core_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gojek/optimus-extension-valor/core"
	_ "github.com/gojek/optimus-extension-valor/plugin/explorer"
	_ "github.com/gojek/optimus-extension-valor/plugin/formatter"
	_ "github.com/gojek/optimus-extension-valor/plugin/io"
	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/stretchr/testify/suite"
)

const (
	defaultDirName            = "./out"
	defaultSchemaFileName     = "schema.json"
	defaultSchemaContent      = "{\"message\":0}"
	defaultDefinitionFileName = "definition.json"
	defaultDefinitionContent  = "{\"message\":0}"
	defaultDefinitionFormat   = "json"
	defaultProcedureFileName  = "procedure.jsonnet"
	defaultProcedureContent   = "{\"message\":0}"
	defaultValidType          = "file"
)

type LoaderSuite struct {
	suite.Suite
}

func (l *LoaderSuite) SetupSuite() {
	if err := os.MkdirAll(defaultDirName, os.ModePerm); err != nil {
		panic(err)
	}
	filePath := path.Join(defaultDirName, defaultSchemaFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultSchemaContent), os.ModePerm); err != nil {
		panic(err)
	}
	filePath = path.Join(defaultDirName, defaultDefinitionFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultDefinitionContent), os.ModePerm); err != nil {
		panic(err)
	}
	filePath = path.Join(defaultDirName, defaultProcedureFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultProcedureContent), os.ModePerm); err != nil {
		panic(err)
	}
}

func (l *LoaderSuite) TestLoadFramework() {
	l.Run("should return nil and error if recipe is nil", func() {
		var rcp *recipe.Framework = nil
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadFramework(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if error during loading definition", func() {
		rcp := &recipe.Framework{
			Name: "test_framework",
			Definitions: []*recipe.Definition{
				{
					Name:   "test_definition",
					Format: defaultDefinitionFormat,
					Type:   defaultValidType,
					Path:   path.Join(defaultDirName, defaultDefinitionFileName),
				},
				nil,
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadFramework(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if error during loading schema", func() {
		rcp := &recipe.Framework{
			Name: "test_framework",
			Definitions: []*recipe.Definition{
				{
					Name:   "test_definition",
					Format: defaultDefinitionFormat,
					Type:   defaultValidType,
					Path:   path.Join(defaultDirName, defaultDefinitionFileName),
				},
			},
			Schemas: []*recipe.Schema{
				{
					Name: "test_schema",
					Type: defaultValidType,
					Path: path.Join(defaultDirName, defaultSchemaFileName),
				},
				nil,
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadFramework(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if error during loading procedure", func() {
		rcp := &recipe.Framework{
			Name: "test_framework",
			Definitions: []*recipe.Definition{
				{
					Name:   "test_definition",
					Format: defaultDefinitionFormat,
					Type:   defaultValidType,
					Path:   path.Join(defaultDirName, defaultDefinitionFileName),
				},
			},
			Schemas: []*recipe.Schema{
				{
					Name: "test_schema",
					Type: defaultValidType,
					Path: path.Join(defaultDirName, defaultSchemaFileName),
				},
			},
			Procedures: []*recipe.Procedure{
				{
					Name: "test_procedure",
					Type: defaultValidType,
					Path: path.Join(defaultDirName, defaultProcedureFileName),
				},
				nil,
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadFramework(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Framework{
			Name: "test_framework",
			Definitions: []*recipe.Definition{
				{
					Name:   "test_definition",
					Format: defaultDefinitionFormat,
					Type:   defaultValidType,
					Path:   path.Join(defaultDirName, defaultDefinitionFileName),
				},
			},
			Schemas: []*recipe.Schema{
				{
					Name: "test_schema",
					Type: defaultValidType,
					Path: path.Join(defaultDirName, defaultSchemaFileName),
				},
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadFramework(rcp)

		l.NotNil(actualValue)
		l.Nil(actualErr)
	})
}

func (l *LoaderSuite) TestLoadProcedure() {
	l.Run("should return nil and error if recipe is nil", func() {
		var rcp *recipe.Procedure = nil
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid type", func() {
		rcp := &recipe.Procedure{
			Name: "test_procedure",
			Type: "invalid_type",
			Path: path.Join(defaultDirName, defaultProcedureFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Procedure{
			Name: "test_procedure",
			Type: defaultValidType,
			Path: path.Join(defaultDirName, defaultProcedureFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.NotNil(actualValue)
		l.Nil(actualErr)
	})
}

func (l *LoaderSuite) TestLoadSchema() {
	l.Run("should return nil and error if recipe is nil", func() {
		var rcp *recipe.Schema = nil
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid type", func() {
		rcp := &recipe.Schema{
			Name: "test_schema",
			Type: "invalid_type",
			Path: path.Join(defaultDirName, defaultSchemaFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Schema{
			Name: "test_schema",
			Type: defaultValidType,
			Path: path.Join(defaultDirName, defaultSchemaFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

		l.NotNil(actualValue)
		l.Nil(actualErr)
	})
}

func (l *LoaderSuite) TestLoadDefinition() {
	l.Run("should return nil and error if recipe is nil", func() {
		var rcp *recipe.Definition = nil
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid type", func() {
		rcp := &recipe.Definition{
			Name:   "test_definition",
			Type:   "invalid_type",
			Format: defaultDefinitionFormat,
			Path:   path.Join(defaultDirName, defaultDefinitionFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe function is set but empty", func() {
		rcp := &recipe.Definition{
			Name:     "test_definition",
			Type:     defaultValidType,
			Format:   defaultDefinitionFormat,
			Path:     path.Join(defaultDirName, defaultDefinitionFileName),
			Function: &recipe.Function{},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Definition{
			Name:   "test_definition",
			Type:   defaultValidType,
			Format: defaultDefinitionFormat,
			Path:   path.Join(defaultDirName, defaultDefinitionFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.NotNil(actualValue)
		l.Nil(actualErr)
	})
}

func (l *LoaderSuite) TestLoadData() {
	l.Run("should return nil and error if type is invalid", func() {
		loader := &core.Loader{}
		pt := path.Join(defaultDirName, defaultDefinitionFileName)
		_type := "invalid_type"
		format := defaultDefinitionFormat

		actualData, actualErr := loader.LoadData(pt, _type, format)

		l.Nil(actualData)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if format is invalid", func() {
		loader := &core.Loader{}
		pt := path.Join(defaultDirName, defaultDefinitionFileName)
		_type := defaultValidType
		format := "invalid_format"

		actualData, actualErr := loader.LoadData(pt, _type, format)

		l.Nil(actualData)
		l.NotNil(actualErr)
	})

	l.Run("should return data and nil if no error is encountered", func() {
		loader := &core.Loader{}
		pt := path.Join(defaultDirName, defaultDefinitionFileName)
		_type := defaultValidType
		format := defaultDefinitionFormat

		actualData, actualErr := loader.LoadData(pt, _type, format)

		l.NotNil(actualData)
		l.Nil(actualErr)
	})
}

func (l *LoaderSuite) TearDownSuite() {
	if err := os.RemoveAll(defaultDirName); err != nil {
		panic(err)
	}
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, &LoaderSuite{})
}
