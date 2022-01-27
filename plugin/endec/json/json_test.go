package json_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/plugin/endec/json"

	"github.com/stretchr/testify/assert"
)

func TestNewEncode(t *testing.T) {
	t.Run("should return nil and error if error when executing", func(t *testing.T) {
		message := func() {}
		encode := json.NewEncode()

		actualContent, actualErr := encode(message)

		assert.Nil(t, actualContent)
		assert.NotNil(t, actualErr)
	})

	t.Run("should return content and nil if no error is encountered during execution", func(t *testing.T) {
		message := "message"
		encode := json.NewEncode()

		actualContent, actualErr := encode(message)

		assert.NotNil(t, actualContent)
		assert.Nil(t, actualErr)
	})
}

func TestNewDecode(t *testing.T) {
	t.Run("should return error if error when executing", func(t *testing.T) {
		input := []byte("message")
		var output int
		decode := json.NewDecode()

		actualErr := decode(input, &output)

		assert.NotNil(t, actualErr)
	})

	t.Run("should return content and nil if no error is encountered during execution", func(t *testing.T) {
		input := []byte("{\"message\": 0}")
		var output interface{}
		decode := json.NewDecode()

		actualErr := decode(input, &output)

		assert.Nil(t, actualErr)
	})
}
