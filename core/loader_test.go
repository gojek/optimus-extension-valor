package core_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gojek/optimus-extension-valor/core"
	_ "github.com/gojek/optimus-extension-valor/plugin/formatter"
	_ "github.com/gojek/optimus-extension-valor/plugin/io"
	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/stretchr/testify/suite"
)

const (
	defaulDirName          = "./out"
	defaultValidFileName   = "valor.yaml"
	defaultInvalidFileName = "valor.invalid"
	defaultContent         = "message: 0"
)

type LoaderSuite struct {
	suite.Suite
}

func (l *LoaderSuite) SetupSuite() {
	if err := os.MkdirAll(defaulDirName, os.ModePerm); err != nil {
		panic(err)
	}
	filePath := path.Join(defaulDirName, defaultValidFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultContent), os.ModePerm); err != nil {
		panic(err)
	}
	filePath = path.Join(defaulDirName, defaultInvalidFileName)
	if err := ioutil.WriteFile(filePath, []byte(defaultContent), os.ModePerm); err != nil {
		panic(err)
	}
}

func (l *LoaderSuite) TestLoadResource() {
	l.Run("should return nil and error if recipe is nil", func() {
		var rcp *recipe.Resource = nil
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadResource(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid type", func() {
		rcp := &recipe.Resource{
			Name:   "test_resource",
			Type:   "invalid_type",
			Format: "yaml",
			Path:   path.Join(defaulDirName, defaultValidFileName),
			FrameworkNames: []string{
				"framework_target",
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadResource(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid format", func() {
		rcp := &recipe.Resource{
			Name:   "test_resource",
			Type:   "file",
			Format: "invalid_format",
			Path:   path.Join(defaulDirName, defaultValidFileName),
			FrameworkNames: []string{
				"framework_target",
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadResource(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains inconsistent format", func() {
		rcp := &recipe.Resource{
			Name:   "test_resource",
			Type:   "dir",
			Format: "inconsistent",
			Path:   defaulDirName,
			FrameworkNames: []string{
				"framework_target",
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadResource(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid path", func() {
		rcp := &recipe.Resource{
			Name:   "test_resource",
			Type:   "file",
			Format: "yaml",
			Path:   defaulDirName,
			FrameworkNames: []string{
				"framework_target",
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadResource(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Resource{
			Name:   "test_resource",
			Type:   "file",
			Format: "yaml",
			Path:   path.Join(defaulDirName, defaultValidFileName),
			FrameworkNames: []string{
				"framework_target",
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadResource(rcp)

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
			Format: "yaml",
			Path:   path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid format", func() {
		rcp := &recipe.Definition{
			Name:   "test_definition",
			Type:   "file",
			Format: "invalid_format",
			Path:   path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains inconsistent format", func() {
		rcp := &recipe.Definition{
			Name:   "test_definition",
			Type:   "file",
			Format: "inconsistent",
			Path:   path.Join(defaulDirName, defaultInvalidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid path", func() {
		rcp := &recipe.Definition{
			Name:   "test_definition",
			Type:   "file",
			Format: "yaml",
			Path:   defaulDirName,
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe function is empty", func() {
		rcp := &recipe.Definition{
			Name:     "test_definition",
			Type:     "file",
			Format:   "yaml",
			Path:     path.Join(defaulDirName, defaultValidFileName),
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
			Type:   "file",
			Format: "yaml",
			Path:   path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadDefinition(rcp)

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
			Path: path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains inconsistent format", func() {
		rcp := &recipe.Schema{
			Name: "test_schema",
			Type: "dir",
			Path: defaulDirName,
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid path", func() {
		rcp := &recipe.Schema{
			Name: "test_schema",
			Type: "file",
			Path: defaulDirName,
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Schema{
			Name: "test_schema",
			Type: "file",
			Path: path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadSchema(rcp)

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
			Path: path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains inconsistent format", func() {
		rcp := &recipe.Procedure{
			Name: "test_procedure",
			Type: "dir",
			Path: defaulDirName,
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return nil and error if recipe contains invalid path", func() {
		rcp := &recipe.Procedure{
			Name: "test_procedure",
			Type: "file",
			Path: defaulDirName,
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.Nil(actualValue)
		l.NotNil(actualErr)
	})

	l.Run("should return value and nil if no error is encountered", func() {
		rcp := &recipe.Procedure{
			Name: "test_procedure",
			Type: "file",
			Path: path.Join(defaulDirName, defaultValidFileName),
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadProcedure(rcp)

		l.NotNil(actualValue)
		l.Nil(actualErr)
	})
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
					Format: "yaml",
					Type:   "file",
					Path:   path.Join(defaulDirName, defaultValidFileName),
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
					Format: "yaml",
					Type:   "file",
					Path:   path.Join(defaulDirName, defaultValidFileName),
				},
			},
			Schemas: []*recipe.Schema{
				{
					Name: "test_schema",
					Type: "file",
					Path: path.Join(defaulDirName, defaultValidFileName),
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
					Format: "yaml",
					Type:   "file",
					Path:   path.Join(defaulDirName, defaultValidFileName),
				},
			},
			Schemas: []*recipe.Schema{
				{
					Name: "test_schema",
					Type: "file",
					Path: path.Join(defaulDirName, defaultValidFileName),
				},
			},
			Procedures: []*recipe.Procedure{
				{
					Name: "test_procedure",
					Type: "file",
					Path: path.Join(defaulDirName, defaultValidFileName),
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
					Format: "yaml",
					Type:   "file",
					Path:   path.Join(defaulDirName, defaultValidFileName),
				},
			},
			Schemas: []*recipe.Schema{
				{
					Name: "test_schema",
					Type: "file",
					Path: path.Join(defaulDirName, defaultValidFileName),
				},
			},
			Procedures: []*recipe.Procedure{
				{
					Name: "test_procedure",
					Type: "file",
					Path: path.Join(defaulDirName, defaultValidFileName),
				},
			},
			OutputTargets: []*recipe.OutputTarget{
				{
					Name:   "std",
					Format: "yaml",
					Type:   "file",
					Path:   path.Join(defaulDirName),
				},
			},
		}
		loader := &core.Loader{}

		actualValue, actualErr := loader.LoadFramework(rcp)

		l.NotNil(actualValue)
		l.Nil(actualErr)
	})
}

func (l *LoaderSuite) TearDownSuite() {
	if err := os.RemoveAll(defaulDirName); err != nil {
		panic(err)
	}
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, &LoaderSuite{})
}
