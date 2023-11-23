package core

import (
	"regexp"
	"strings"

	"github.com/gojek/optimus-extension-valor/registry/explorer"
)

// ExplorePaths explores the given root path for the type and format
func ExplorePaths(rootPath, _type, format, regexPattern string) ([]string, error) {
	exPath, err := explorer.Explorers.Get(_type)
	if err != nil {
		return nil, err
	}
	reg, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}
	return exPath(rootPath, func(path string) bool {
		return reg.MatchString(path) && strings.HasSuffix(path, format)
	})
}
