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
	t.Run("should return nil and error if reader is nil", func(t *testing.T) {
		var reader model.Reader = nil
		var decode model.Decode

		expectedErr := errors.New("reader is nil")

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if decode is nil", func(t *testing.T) {
		reader := &mocks.Reader{}
		var decode model.Decode = nil

		expectedErr := errors.New("decode is nil")

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if read returns error", func(t *testing.T) {
		readErr := errors.New("read error")
		reader := &mocks.Reader{}
		reader.On("Read").Return(nil, readErr)
		decode := func(c []byte, v interface{}) error {
			return nil
		}

		expectedErr := readErr

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if content is empty", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("Read").Return(&model.Data{}, nil)
		decode := func(c []byte, v interface{}) error {
			return nil
		}

		expectedErr := errors.New("content is empty")

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return nil and error if decode returns error", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("Read").Return(&model.Data{
			Content: []byte("message"),
		}, nil)
		decodeErr := errors.New("decode error")
		decode := func(c []byte, v interface{}) error {
			return decodeErr
		}

		expectedErr := decodeErr

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.Nil(t, actualRecipe)
		assert.EqualValues(t, expectedErr, actualErr)
	})

	t.Run("should return data and nil if no error is encountered", func(t *testing.T) {
		reader := &mocks.Reader{}
		reader.On("Read").Return(&model.Data{
			Content: []byte("message"),
		}, nil)
		decode := func(c []byte, v interface{}) error {
			return nil
		}

		actualRecipe, actualErr := recipe.Load(reader, decode)

		assert.NotNil(t, actualRecipe)
		assert.Nil(t, actualErr)
	})
}
