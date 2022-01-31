package core_test

import (
	"errors"
	"testing"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/model"
	_ "github.com/gojek/optimus-extension-valor/plugin/endec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EvaluatorSuite struct {
	suite.Suite
}

func (e *EvaluatorSuite) TestEvaluate() {
	e.Run("should return false and error if resource data is nil", func() {
		framework := &model.Framework{}
		var resourceData *model.Data = nil
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		evaluator, _ := core.NewEvaluator(framework, evaluate)

		expectedValue := false

		actualValue, actualErr := evaluator.Evaluate(resourceData)

		e.Equal(expectedValue, actualValue)
		e.NotNil(actualErr)
	})

	e.Run("should return false and error if one or more procedures are nil", func() {
		framework := &model.Framework{
			Procedures: []*model.Procedure{nil},
		}
		var resourceData *model.Data = &model.Data{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		evaluator, _ := core.NewEvaluator(framework, evaluate)

		expectedValue := false

		actualValue, actualErr := evaluator.Evaluate(resourceData)

		e.Equal(expectedValue, actualValue)
		e.NotNil(actualErr)
	})

	e.Run("should return null and error if procedure data is nil", func() {
		framework := &model.Framework{
			Procedures: []*model.Procedure{
				{
					Name: "procecure_test",
				},
			},
		}
		var resourceData *model.Data = &model.Data{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		evaluator, _ := core.NewEvaluator(framework, evaluate)

		expectedValue := false

		actualValue, actualErr := evaluator.Evaluate(resourceData)

		e.Equal(expectedValue, actualValue)
		e.NotNil(actualErr)
	})

	e.Run("should return null and error if evaluation results in error", func() {
		framework := &model.Framework{
			Procedures: []*model.Procedure{
				{
					Name: "procecure_test",
					Data: &model.Data{
						Content: []byte("test content"),
					},
				},
			},
		}
		var resourceData *model.Data = &model.Data{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", errors.New("test error")
		}
		evaluator, _ := core.NewEvaluator(framework, evaluate)

		expectedValue := false

		actualValue, actualErr := evaluator.Evaluate(resourceData)

		e.Equal(expectedValue, actualValue)
		e.NotNil(actualErr)
	})

	e.Run("should return true and nil if no error is encountered", func() {
		framework := &model.Framework{
			Procedures: []*model.Procedure{
				{
					Name: "procecure_test",
					Data: &model.Data{
						Content: []byte("test content"),
					},
				},
			},
		}
		var resourceData *model.Data = &model.Data{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "{\"message\": \"error\"}", nil
		}
		evaluator, _ := core.NewEvaluator(framework, evaluate)

		expectedValue := true

		actualValue, actualErr := evaluator.Evaluate(resourceData)

		e.Equal(expectedValue, actualValue)
		e.Nil(actualErr)
	})
}

func TestEvaluatorSuite(t *testing.T) {
	suite.Run(t, &EvaluatorSuite{})
}

func TestNewEvaluator(t *testing.T) {
	t.Run("should return nil and error if framework is nil", func(t *testing.T) {
		var framework *model.Framework = nil
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}

		actualValue, actualErr := core.NewEvaluator(framework, evaluate)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if evaluate is nil", func(t *testing.T) {
		framework := &model.Framework{}
		var evaluate model.Evaluate = nil

		actualValue, actualErr := core.NewEvaluator(framework, evaluate)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if one or more definition is nil", func(t *testing.T) {
		framework := &model.Framework{
			Definitions: []*model.Definition{
				{
					Name: "test_definition",
				},
				nil,
			},
		}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}

		actualValue, actualErr := core.NewEvaluator(framework, evaluate)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if definition contains invalid data", func(t *testing.T) {
		framework := &model.Framework{
			Definitions: []*model.Definition{
				{
					Name: "test_definition",
					ListOfData: []*model.Data{
						{
							Content: []byte("test content"),
						},
					},
					FunctionData: &model.Data{
						Content: []byte("test function"),
					},
				},
			},
		}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", errors.New("test error")
		}

		actualValue, actualErr := core.NewEvaluator(framework, evaluate)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return value and nil if no error is encountered", func(t *testing.T) {
		framework := &model.Framework{
			Definitions: []*model.Definition{
				{
					Name: "test_definition",
					ListOfData: []*model.Data{
						{
							Type:    "file",
							Path:    "test.yaml",
							Content: []byte("message: name"),
						},
					},
				},
			},
		}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}

		actualValue, actualErr := core.NewEvaluator(framework, evaluate)

		assert.NotNil(t, actualValue)
		assert.Nil(t, actualErr)
	})
}
