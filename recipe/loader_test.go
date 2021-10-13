package recipe_test

import (
	"errors"
	"testing"

	"github.com/gojek/optimus-extension-valor/mocks"
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin"
	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/stretchr/testify/assert"
)

func TestLoadWithReader(t *testing.T) {
	t.Run("should return nil and error if reader is nil", func(t *testing.T) {
		var reader plugin.Reader = nil
		decoder := func([]byte, interface{}) error {
			return nil
		}

		expectedErrorMsg := "reader is nil"

		actualRecipe, actualError := recipe.LoadWithReader(reader, decoder)

		assert.Nil(t, actualRecipe)
		assert.EqualError(t, actualError, expectedErrorMsg)
	})

	t.Run("should return nil and error if decoder is nil", func(t *testing.T) {
		reader := &mocks.Reader{}
		var decoder plugin.Decoder = nil

		expectedErrorMsg := "decoder is nil"

		actualRecipe, actualError := recipe.LoadWithReader(reader, decoder)

		assert.Nil(t, actualRecipe)
		assert.EqualError(t, actualError, expectedErrorMsg)
	})

	t.Run("should return nil and error if reader returns error", func(t *testing.T) {
		errorMsg := "cannot read data"
		reader := &mocks.Reader{}
		reader.On("Read").Return(nil, errors.New(errorMsg))
		decoder := func([]byte, interface{}) error {
			return nil
		}

		expectedErrorMsg := "cannot read data"

		actualRecipe, actualError := recipe.LoadWithReader(reader, decoder)

		assert.Nil(t, actualRecipe)
		assert.EqualError(t, actualError, expectedErrorMsg)
	})

	t.Run("should return nil and error if content is empty", func(t *testing.T) {
		content := ""
		reader := &mocks.Reader{}
		reader.On("Read").Return(&model.Data{Content: []byte(content)}, nil)
		decoder := func([]byte, interface{}) error {
			return nil
		}

		expectedErrorMsg := "content is empty"

		actualRecipe, actualError := recipe.LoadWithReader(reader, decoder)

		assert.Nil(t, actualRecipe)
		assert.EqualError(t, actualError, expectedErrorMsg)
	})

	t.Run("should return nil and error if decoder returns error", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("Read").Return(&model.Data{Content: []byte("invalid")}, nil)
		decoder := func([]byte, interface{}) error {
			return errors.New("cannot be decoded")
		}

		expectedErrorMsg := "cannot be decoded"

		actualRecipe, actualError := recipe.LoadWithReader(reader, decoder)

		assert.Nil(t, actualRecipe)
		assert.EqualError(t, actualError, expectedErrorMsg)
	})

	t.Run("should return recipe and nil if no error encountered", func(t *testing.T) {
		content := "key: value"
		reader := &mocks.Reader{}
		reader.On("Read").Return(&model.Data{Content: []byte(content)}, nil)
		decoder := func([]byte, interface{}) error {
			return nil
		}

		actualRecipe, actualError := recipe.LoadWithReader(reader, decoder)

		assert.NotNil(t, actualRecipe)
		assert.Nil(t, actualError)
	})
}
