package recipe_test

import (
	"errors"
	"testing"

	"github.com/gojek/optimus-extension-valor/mocks"
	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	const defaultErrKey = "Load"

	t.Run("should return nil and error if reader is nil", func(t *testing.T) {
		var reader model.Reader = nil
		var decode model.Decode

		expectedErr := model.BuildError(defaultErrKey, errors.New("reader is nil"))

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if decode is nil", func(t *testing.T) {
		reader := &mocks.Reader{}
		var decode model.Decode = nil

		expectedErr := model.BuildError(defaultErrKey, errors.New("decode is nil"))

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if read returns error", func(t *testing.T) {
		readErr := model.BuildError(defaultErrKey, errors.New("read error"))
		reader := &mocks.Reader{}
		reader.On("ReadOne").Return(nil, readErr)
		decode := func(c []byte, v interface{}) model.Error {
			return nil
		}

		expectedErr := readErr

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if content is empty", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("ReadOne").Return(&model.Data{}, nil)
		decode := func(c []byte, v interface{}) model.Error {
			return nil
		}

		expectedErr := model.BuildError(defaultErrKey, errors.New("content is empty"))

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if decode returns error", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("ReadOne").Return(&model.Data{
			Content: []byte("message"),
		}, nil)
		decodeErr := model.BuildError(defaultErrKey, errors.New("decode error"))
		decode := func(c []byte, v interface{}) model.Error {
			return decodeErr
		}

		expectedErr := decodeErr

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return data and nil if no error is encountered", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("ReadOne").Return(&model.Data{
			Content: []byte("message"),
		}, nil)
		decode := func(c []byte, v interface{}) model.Error {
			return nil
		}

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.NotNil(t, actualRecipe)
		assert.Nil(t, actualErr)
	})
}
