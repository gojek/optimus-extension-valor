package json_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/plugin/formatter/json"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	t.Run("should return the same value but not the same address and nil", func(t *testing.T) {
		input := []byte("message")

		expectedValue := input

		actualValue, actualErr := json.ToJSON(input)

		assert.EqualValues(t, expectedValue, actualValue)
		assert.Equal(t, expectedValue, actualValue)
		assert.Nil(t, actualErr)
	})
}

func TestToYAML(t *testing.T) {
	t.Run("should return nil and error if error during unmarshal", func(t *testing.T) {
		input := []byte("message")

		actualValue, actualErr := json.ToYAML(input)

		assert.Nil(t, actualValue)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return value and nil if no error is encountered", func(t *testing.T) {
		input := []byte("{\"message\":0}")

		actualValue, actualErr := json.ToYAML(input)

		assert.NotNil(t, actualValue)
		assert.Nil(t, actualErr)
	})
}
