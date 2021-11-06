package std

import (
	"fmt"
	"os"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "std"

// Std is an stdout writer
type Std struct {
}

func (s *Std) Write(dataList ...*model.Data) model.Error {
	const defaultErrKey = "Write"
	for _, d := range dataList {
		output := fmt.Sprintf("%s\n%s\n", d.Path, string(d.Content))
		_, err := os.Stdout.WriteString(output)
		if err != nil {
			return model.BuildError(defaultErrKey, err)
		}
	}
	return nil
}

// New initializes standard input and output
func New() *Std {
	return &Std{}
}

func init() {
	err := io.Writers.Register(_type, New())
	if err != nil {
		panic(err)
	}
}
