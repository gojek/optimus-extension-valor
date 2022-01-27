package yaml_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/plugin/endec/yaml"
	"github.com/stretchr/testify/assert"
)

func TestNewEncode(t *testing.T) {
	t.Run("should return nil and error if error when executing", func(t *testing.T) {
		message := func() {}
		encode := yaml.NewEncode()

		actualContent, actualErr := encode(message)

		assert.Nil(t, actualContent)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return content and nil if no error is encountered during execution", func(t *testing.T) {
		message := "message"
		encode := yaml.NewEncode()

		actualContent, actualErr := encode(message)

		assert.NotNil(t, actualContent)
		assert.Nil(t, actualErr)
	})
}

func TestNewDecode(t *testing.T) {
	t.Run("should return error if error when executing", func(t *testing.T) {
		input := []byte("message")
		var output int
		decode := yaml.NewDecode()

		actualErr := decode(input, &output)

		assert.NotNil(t, actualErr)
	})

	t.Run("should return content and nil if no error is encountered during execution", func(t *testing.T) {
		input := []byte("message: 0")
		var output interface{}
		decode := yaml.NewDecode()

		actualErr := decode(input, &output)

		assert.Nil(t, actualErr)
	})
}
