package core

import (
	"github.com/google/go-jsonnet"
)

// Evaluate evaluates snippet
func Evaluate(name, snippet string) ([]byte, error) {
	vm := jsonnet.MakeVM()
	rst, err := vm.EvaluateAnonymousSnippet(name, snippet)
	if err != nil {
		return nil, err
	}
	return []byte(rst), nil
}
