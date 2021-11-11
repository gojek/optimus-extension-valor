package yaml_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/plugin/formatter/yaml"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	t.Run("should return value and nil if no error encountered", func(t *testing.T) {
		input := []byte("{\"message\": 0}")

		actualValue, actualErr := yaml.ToJSON(input)

		assert.NotNil(t, actualValue)
		assert.Nil(t, actualErr)
	})
}
