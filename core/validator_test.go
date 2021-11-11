package core_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ValidatorSuite struct {
	suite.Suite
}

func (v *ValidatorSuite) TestValidate() {
	v.Run("should return error if resource data is nil", func() {
		framework := &model.Framework{}
		var resourceData *model.Data = nil
		validator, _ := core.NewValidator(framework)

		actualErr := validator.Validate(resourceData)

		v.NotNil(actualErr)
	})

	v.Run("should return error if schema data is nil", func() {
		framework := &model.Framework{
			Schemas: []*model.Schema{
				{
					Name: "schema_test",
					Data: nil,
				},
			},
		}
		resourceData := &model.Data{}
		validator, _ := core.NewValidator(framework)

		actualErr := validator.Validate(resourceData)

		v.NotNil(actualErr)
	})

	v.Run("should return error if validation returns error", func() {
		schemaContent := `{
    "title": "user_account",
    "description": "Schema to validate user_account.",
    "type": "object",
    "properties": {
        "email": {
            "type": "string"
        },
        "membership": {
            "enum": [
                "standard",
                "premium"
            ]
        },
        "is_active": {
            "type": "boolean"
        }
    },
    "required": [
		"email",
		"membership"
	],
    "additionalProperties": false
}
`
		framework := &model.Framework{
			Schemas: []*model.Schema{
				{
					Name: "schema_test",
					Data: &model.Data{
						Content: []byte(schemaContent),
					},
				},
			},
		}
		resourceContent := `{
    "email": "valor@github.com",
    "membership": "premium",
    "is_active": 1
}
`
		resourceData := &model.Data{
			Content: []byte(resourceContent),
		}
		validator, _ := core.NewValidator(framework)

		actualErr := validator.Validate(resourceData)

		v.NotNil(actualErr)
	})

	v.Run("should return nil if validation success", func() {
		schemaContent := `{
    "title": "user_account",
    "description": "Schema to validate user_account.",
    "type": "object",
    "properties": {
        "email": {
            "type": "string"
        },
        "membership": {
            "enum": [
                "standard",
                "premium"
            ]
        },
        "is_active": {
            "type": "boolean"
        }
    },
    "required": [
		"email",
		"membership"
	],
    "additionalProperties": false
}
`
		framework := &model.Framework{
			Schemas: []*model.Schema{
				{
					Name: "schema_test",
					Data: &model.Data{
						Content: []byte(schemaContent),
					},
				},
			},
		}
		resourceContent := `{
    "email": "valor@github.com",
    "membership": "premium",
    "is_active": true
}
`
		resourceData := &model.Data{
			Content: []byte(resourceContent),
		}
		validator, _ := core.NewValidator(framework)

		actualErr := validator.Validate(resourceData)

		v.Nil(actualErr)
	})
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, &ValidatorSuite{})
}

func TestNewValidator(t *testing.T) {
	t.Run("should return nil and error if framework is nil", func(t *testing.T) {
		var framework *model.Framework = nil

		actualValue, actualErr := core.NewValidator(framework)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if one or more schema is nil", func(t *testing.T) {
		framework := &model.Framework{
			Schemas: []*model.Schema{
				{
					Name: "test_schema",
				},
				nil,
			},
		}

		actualValue, actualErr := core.NewValidator(framework)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return validator and nil if no error is encountered", func(t *testing.T) {
		framework := &model.Framework{
			Schemas: []*model.Schema{
				{
					Name: "test_schema",
				},
			},
		}

		actualValue, actualErr := core.NewValidator(framework)

		assert.NotNil(t, actualValue)
		assert.Nil(t, actualErr)
	})
}
