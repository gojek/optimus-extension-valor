package core

import (
	"path"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

func treatOutput(data *model.Data, output *model.Output) (bool, error) {
	if output == nil {
		return true, nil
	}
	outputError := make(model.Error)
	for _, t := range output.Targets {
		formatterFn, err := formatter.Formats.Get(jsonFormat, t.Format)
		if err != nil {
			outputError[t.Name] = err
			continue
		}
		writerFn, err := io.Writers.Get(t.Type)
		if err != nil {
			outputError[t.Name] = err
			continue
		}
		result, err := formatterFn(data.Content)
		if err != nil {
			outputError[t.Name] = err
			continue
		}
		writer := writerFn(output.TreatAs)
		if err := writer.Write(
			&model.Data{
				Type:    data.Type,
				Path:    path.Join(t.Path, data.Path),
				Content: result,
			},
		); err != nil {
			outputError[t.Name] = err
			continue
		}
	}
	if len(outputError) > 0 {
		return false, outputError
	}
	if len(output.Targets) > 0 && output.TreatAs == model.TreatmentError {
		return false, nil
	}
	return true, nil
}
