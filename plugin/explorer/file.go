package explorer

import (
	"io/fs"
	"path/filepath"

	"github.com/gojek/optimus-extension-valor/registry/explorer"
)

const fileType = "file"

// ExploreFilePath explores file path from its root based on its filter
func ExploreFilePath(root string, filter func(string) bool) ([]string, error) {
	var output []string
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if filter == nil || filter(path) {
				output = append(output, path)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return output, nil
}

func init() {
	if err := explorer.Explorers.Register(fileType, ExploreFilePath); err != nil {
		panic(err)
	}
}
