package std

import (
	"os"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/fatih/color"
)

const _type = "std"

// Std is an stdout writer
type Std struct {
	treatment model.OutputTreatment
}

func (s *Std) Write(dataList ...*model.Data) model.Error {
	const defaultErrKey = "Write"
	color.NoColor = false
	var colorize *color.Color
	switch s.treatment {
	case model.TreatmentError:
		colorize = color.New(color.FgHiRed)
	case model.TreatmentWarning:
		colorize = color.New(color.FgHiYellow)
	case model.TreatmentSuccess:
		colorize = color.New(color.FgHiGreen)
	default:
		colorize = color.New(color.FgHiWhite)
	}
	for _, d := range dataList {
		separator := strings.Repeat("-", len(d.Path))
		output := colorize.Sprintf("%s\n%s\n%s\n%s\n", separator, d.Path, separator, string(d.Content))
		_, err := os.Stdout.WriteString(output)
		if err != nil {
			return model.BuildError(defaultErrKey, err)
		}
	}
	return nil
}

// New initializes standard input and output
func New(treatment model.OutputTreatment) *Std {
	return &Std{
		treatment: treatment,
	}
}

func init() {
	err := io.Writers.Register(_type, func(treatment model.OutputTreatment) model.Writer {
		return New(treatment)
	})
	if err != nil {
		panic(err)
	}
}
