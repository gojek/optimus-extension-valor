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
	v.Run("should return false and error if resource data is nil", func() {
		framework := &model.Framework{}
		var resourceData *model.Data = nil
		validator, _ := core.NewValidator(framework)

		actualSuccess, actualErr := validator.Validate(resourceData)

		v.False(actualSuccess)
		v.NotNil(actualErr)
	})

	v.Run("should return false and error if schema data is nil", func() {
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

		actualSuccess, actualErr := validator.Validate(resourceData)

		v.False(actualSuccess)
		v.NotNil(actualErr)
	})

	v.Run("should return false and nil if validation execution success but is business error", func() {
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
    "membership": "invalid",
    "is_active": true
}
`
		resourceData := &model.Data{
			Content: []byte(resourceContent),
		}
		validator, _ := core.NewValidator(framework)

		actualSuccess, actualErr := validator.Validate(resourceData)

		v.True(actualSuccess)
		v.Nil(actualErr)
	})

	v.Run("should return true and nil if validation execution and business are success", func() {
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

		actualSuccess, actualErr := validator.Validate(resourceData)

		v.True(actualSuccess)
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
