package std

import (
	"fmt"
	"os"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/plugin"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const _type = "std"

// Std is an stdout writer
type Std struct {
}

func (s *Std) Write(dataList ...*model.Data) error {
	const nameKey = "name"
	const pathKey = "path"
	for _, d := range dataList {
		var missingKeys []string
		if d.Metadata[nameKey] == "" {
			missingKeys = append(missingKeys, nameKey)
		}
		if d.Metadata[pathKey] == "" {
			missingKeys = append(missingKeys, pathKey)
		}
		if len(missingKeys) > 0 {
			return fmt.Errorf("[%s] are empty in metadata", strings.Join(missingKeys, ", "))
		}
		output := fmt.Sprintf("%s: %s\n%s\n", d.Metadata[nameKey], d.Path, string(d.Content))
		_, err := os.Stdout.WriteString(output)
		if err != nil {
			return err
		}
	}
	return nil
}

// New initializes standard input and output
func New() *Std {
	return &Std{}
}

func init() {
	io.Writers.Register(_type, func(path string, metadata map[string]string) plugin.Writer {
		return New()
	})
}
