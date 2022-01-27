package std

import (
	"errors"
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

func (s *Std) Write(data *model.Data) error {
	if data == nil {
		return errors.New("data is nil")
	}
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
	separator := strings.Repeat("-", len(data.Path))
	output := colorize.Sprintf("%s\n%s\n%s\n%s\n",
		separator,
		data.Path,
		separator,
		string(data.Content),
	)
	_, err := os.Stdout.WriteString(output)
	return err
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
