package recipe_test

import (
	"errors"
	"testing"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	const defaultErrKey = "Validate"

	t.Run("should return error if validator returns error", func(t *testing.T) {
		var rcp *recipe.Recipe = nil

		actualErr := recipe.Validate(rcp)

		assert.Error(t, actualErr)
	})

	t.Run("should return error if one or more recipes are invalid", func(t *testing.T) {
		rcp := &recipe.Recipe{
			Resources: []*recipe.Resource{
				{
					Name:           "resource1",
					Format:         "yaml",
					Type:           "file",
					Path:           "./valor.yaml",
					FrameworkNames: []string{"evaluate"},
				},
				{
					Name:   "resource1",
					Format: "yaml",
					Type:   "file",
				},
			},
			Frameworks: []*recipe.Framework{
				{
					Name: "evaluate",
				},
			},
		}

		actualErr := recipe.Validate(rcp)

		assert.NotNil(t, actualErr)
	})

	t.Run("should return error if recipe has duplicate resources", func(t *testing.T) {
		rcp := &recipe.Recipe{
			Resources: []*recipe.Resource{
				{
					Name:           "resource1",
					Format:         "yaml",
					Type:           "file",
					Path:           "./valor.yaml",
					FrameworkNames: []string{"evaluate"},
				},
				{
					Name:           "resource1",
					Format:         "yaml",
					Type:           "file",
					Path:           "./valor.yaml",
					FrameworkNames: []string{"evaluate"},
				},
			},
			Frameworks: []*recipe.Framework{
				{
					Name: "evaluate",
				},
			},
		}

		expectedErr := model.BuildError(defaultErrKey, errors.New("duplicate resource recipe [resource1]"))

		actualErr := recipe.Validate(rcp)

		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return error if one or more recipes are invalid", func(t *testing.T) {
		rcp := &recipe.Recipe{
			Resources: []*recipe.Resource{
				{
					Name:           "resource",
					Format:         "yaml",
					Type:           "file",
					Path:           "./valor.yaml",
					FrameworkNames: []string{"evaluate"},
				},
			},
			Frameworks: []*recipe.Framework{
				{
					Name: "evaluate1",
				},
				{
					Name: "",
				},
			},
		}

		actualErr := recipe.Validate(rcp)

		assert.NotNil(t, actualErr)
	})

	t.Run("should return error if recipe has duplicate frameworks", func(t *testing.T) {
		rcp := &recipe.Recipe{
			Resources: []*recipe.Resource{
				{
					Name:           "resource",
					Format:         "yaml",
					Type:           "file",
					Path:           "./valor.yaml",
					FrameworkNames: []string{"evaluate"},
				},
			},
			Frameworks: []*recipe.Framework{
				{
					Name: "evaluate1",
				},
				{
					Name: "evaluate1",
				},
			},
		}

		expectedErr := model.BuildError(defaultErrKey, errors.New("duplicate framework recipe [evaluate1]"))

		actualErr := recipe.Validate(rcp)

		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil if no error is encountered", func(t *testing.T) {
		rcp := &recipe.Recipe{
			Resources: []*recipe.Resource{
				{
					Name:           "resource",
					Format:         "yaml",
					Type:           "file",
					Path:           "./valor.yaml",
					FrameworkNames: []string{"evaluate"},
				},
			},
			Frameworks: []*recipe.Framework{
				{
					Name: "evaluate1",
				},
			},
		}

		actualErr := recipe.Validate(rcp)

		assert.Nil(t, actualErr)
	})
}

func TestValidateResource(t *testing.T) {
	const defaultErrKey = "ValidateResource"

	t.Run("should return error if validator returns error", func(t *testing.T) {
		var rcp *recipe.Resource = nil

		actualErr := recipe.ValidateResource(rcp)

		assert.Error(t, actualErr)
	})

	t.Run("should return nil if validator returns nil", func(t *testing.T) {
		rcp := &recipe.Resource{
			Name:           "resource1",
			Format:         "yaml",
			Type:           "file",
			Path:           "./valor.yaml",
			FrameworkNames: []string{"evaluate"},
		}

		actualErr := recipe.ValidateResource(rcp)

		assert.Nil(t, actualErr)
	})
}

func TestValidateFramework(t *testing.T) {
	const defaultErrKey = "ValidateFramework"

	t.Run("should return error if validator returns error", func(t *testing.T) {
		var rcp *recipe.Framework = nil

		actualErr := recipe.ValidateFramework(rcp)

		assert.Error(t, actualErr)
	})

	t.Run("should return nil if validator returns nil", func(t *testing.T) {
		rcp := &recipe.Framework{
			Name: "resource1",
		}

		actualErr := recipe.ValidateFramework(rcp)

		assert.Nil(t, actualErr)
	})
}
