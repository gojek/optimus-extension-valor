package core_test

import (
	"testing"

	"github.com/gojek/optimus-extension-valor/core"
	"github.com/gojek/optimus-extension-valor/registry/explorer"

	"github.com/stretchr/testify/assert"
)

func TestExplorePaths(t *testing.T) {
	originalExplorer := explorer.Explorers
	defer func() { explorer.Explorers = originalExplorer }()

	const (
		rootPath = "."
		_type    = "virtual"
		format   = "go"
	)

	t.Run("should return nil and error if explorer registry returns error", func(t *testing.T) {
		_type := "invalid_type"
		pattern := `child_.+`

		actualPaths, actualErr := core.ExplorePaths(rootPath, _type, format, pattern)

		assert.Nil(t, actualPaths)
		assert.Error(t, actualErr)
	})

	explorer.Explorers.Register(_type, func(root string, filter func(string) bool) ([]string, error) {
		paths := []string{
			"./core/testing/root.go",
			"./core/testing/child_1.go",
			"./core/testing/child_2.json",
		}

		var validPaths []string
		for _, p := range paths {
			if filter(p) {
				validPaths = append(validPaths, p)
			}
		}

		return validPaths, nil
	})

	t.Run("should return nil and error if regex pattern is invalid", func(t *testing.T) {
		pattern := "*"

		actualPaths, actualErr := core.ExplorePaths(rootPath, _type, format, pattern)

		assert.Nil(t, actualPaths)
		assert.Error(t, actualErr)
	})

	t.Run("should return all paths and nil if regex pattern is empty", func(t *testing.T) {
		pattern := ""

		actualPaths, actualErr := core.ExplorePaths(rootPath, _type, format, pattern)

		assert.Len(t, actualPaths, 2)
		assert.NoError(t, actualErr)
	})

	t.Run("should return as expected paths and nil if given regex pattern and format", func(t *testing.T) {
		pattern := `child_.+`

		actualPaths, actualErr := core.ExplorePaths(rootPath, _type, format, pattern)

		assert.Len(t, actualPaths, 1)
		assert.NoError(t, actualErr)
	})
}
