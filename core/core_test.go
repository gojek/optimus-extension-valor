package core_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/stretchr/testify/assert"
)

func TestNewPipeline(t *testing.T) {
	t.Run("should return nil and error if recipe is nil", func(t *testing.T) {
		var rcp *recipe.Recipe = nil
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		var batchSize = 0
		var newProgress model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}

		actualPipeline, actualErr := core.NewPipeline(rcp, evaluate, batchSize, newProgress)

		assert.Nil(t, actualPipeline)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if evaluate is nil", func(t *testing.T) {
		var rcp *recipe.Recipe = &recipe.Recipe{}
		var evaluate model.Evaluate = nil
		var batchSize = 0
		var newProgress model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}

		actualPipeline, actualErr := core.NewPipeline(rcp, evaluate, batchSize, newProgress)

		assert.Nil(t, actualPipeline)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if batchSize is less than zero", func(t *testing.T) {
		var rcp *recipe.Recipe = &recipe.Recipe{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		var batchSize = -1
		var newProgress model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}

		actualPipeline, actualErr := core.NewPipeline(rcp, evaluate, batchSize, newProgress)

		assert.Nil(t, actualPipeline)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return nil and error if newProgress is nil", func(t *testing.T) {
		var rcp *recipe.Recipe = &recipe.Recipe{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		var batchSize = 0
		var newProgress model.NewProgress = nil

		actualPipeline, actualErr := core.NewPipeline(rcp, evaluate, batchSize, newProgress)

		assert.Nil(t, actualPipeline)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return pipeline and nil if no error is encountered", func(t *testing.T) {
		var rcp *recipe.Recipe = &recipe.Recipe{}
		var evaluate model.Evaluate = func(name, snippet string) (string, error) {
			return "", nil
		}
		var batchSize = 0
		var newProgress model.NewProgress = func(name string, total int) model.Progress {
			return nil
		}

		actualPipeline, actualErr := core.NewPipeline(rcp, evaluate, batchSize, newProgress)

		assert.NotNil(t, actualPipeline)
		assert.Nil(t, actualErr)
	})
}
