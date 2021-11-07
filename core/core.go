package core

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/google/go-jsonnet"
)

const (
	jsonFormat    = "json"
	jsonnetFormat = "jsonnet"
)

var skipReformat = map[string]bool{
	jsonFormat:    true,
	jsonnetFormat: true,
}

type evaluateWrapper struct {
	Result string
	Error  model.Error
}

// Pipeline defines how a pipeline is executed
type Pipeline struct {
	recipe *recipe.Recipe
	loader *Loader
	vm     *jsonnet.VM

	nameToFrameworkRecipe map[string]*recipe.Framework
}

// NewPipeline initializes pipeline process
func NewPipeline(rcp *recipe.Recipe) (*Pipeline, model.Error) {
	const defaultErrKey = "NewPipeline"
	if rcp == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("recipe is nil"))
	}
	nameToFrameworkRecipe := make(map[string]*recipe.Framework)
	for _, frameworkRcp := range rcp.Frameworks {
		nameToFrameworkRecipe[frameworkRcp.Name] = frameworkRcp
	}
	return &Pipeline{
		recipe:                rcp,
		loader:                &Loader{},
		vm:                    jsonnet.MakeVM(),
		nameToFrameworkRecipe: nameToFrameworkRecipe,
	}, nil
}

// Execute executes pipeline process
func (p *Pipeline) Execute() model.Error {
	const defaultErrKey = "Execute"
	for _, resourceRcp := range p.recipe.Resources {
		decorate := strings.Repeat("=", 12)
		fmt.Printf("%s PROCESSING RESOURCE [%s] %s\n", decorate, resourceRcp.Name, decorate)
		resource, err := p.loader.LoadResource(resourceRcp)
		if err != nil {
			return model.BuildError(defaultErrKey, err)
		}
		for _, frameworkName := range resource.FrameworkNames {
			decoreate := strings.Repeat(":", 5)
			fmt.Printf("%s Processing Framework [%s] %s\n", decoreate, frameworkName, decoreate)
			frameworkRcp := p.nameToFrameworkRecipe[frameworkName]
			fmt.Println(">> Loading framework")
			framework, err := p.loader.LoadFramework(frameworkRcp)
			if err != nil {
				fmt.Println("  Loading failed <<")
				return model.BuildError(defaultErrKey, err)
			}
			fmt.Println(">> Initializing valiator")
			validator, err := NewValidator(framework)
			if err != nil {
				fmt.Println("  Initialization failed <<")
				key := fmt.Sprintf("%s [%s]", defaultErrKey, frameworkName)
				return model.BuildError(key, err)
			}
			fmt.Println(">> Dispatching validation")
			validateChans := p.dispatchValidate(validator, resource.ListOfData)
			outputError := make(model.Error)
			valProgress := NewProgress("getting validation", len(validateChans))
			for i, ch := range validateChans {
				err := <-ch
				if err != nil {
					key := fmt.Sprintf("%s [%s: %s]", defaultErrKey, frameworkName, resource.ListOfData[i].Path)
					outputError[key] = err
				}
				valProgress.Increment()
			}
			valProgress.Wait()
			if len(outputError) > 0 {
				fmt.Println("  Validation failed <<")
				return outputError
			}
			fmt.Println(">> Initializing evaluator")
			evaluator, err := NewEvaluator(framework, p.vm)
			if err != nil {
				fmt.Println("  Initialization failed <<")
				key := fmt.Sprintf("%s [%s]", defaultErrKey, frameworkName)
				return model.BuildError(key, err)
			}
			fmt.Println(">> Dispatching evaluation")
			evalChans := p.dispatchEvaluate(evaluator, resource.ListOfData)
			results := make([]string, len(resource.ListOfData))
			evalProgress := NewProgress("getting evaluation", len(evalChans))
			for i, ch := range evalChans {
				rst := <-ch
				if rst.Error != nil {
					key := fmt.Sprintf("%s [%s: %d]", defaultErrKey, frameworkName, i)
					if resource.ListOfData[i] != nil {
						key = fmt.Sprintf("%s [%s: %s]", defaultErrKey, frameworkName, resource.ListOfData[i].Path)
					}
					outputError[key] = rst.Error
				}
				results[i] = rst.Result
				evalProgress.Increment()
			}
			evalProgress.Wait()
			if len(outputError) > 0 {
				fmt.Println("  Evaluation failed <<")
				return outputError
			}
			fmt.Println(">> Writing output")
			err = p.writeOutput(resource.ListOfData, results, framework.OutputTargets)
			if err != nil {
				fmt.Println("  Write failed <<")
				return model.BuildError(defaultErrKey, err)
			}
		}
		fmt.Println()
	}
	return nil
}

func (p *Pipeline) writeOutput(
	resourceData []*model.Data, resourceResult []string,
	outputTargets []*model.OutputTarget,
) model.Error {
	const defaultErrKey = "writeOutput"
	formatters, err := p.getOutputFormatter(outputTargets)
	if err != nil {
		return model.BuildError(defaultErrKey, err)
	}
	writers, err := p.getOutputWriter(outputTargets)
	if err != nil {
		return model.BuildError(defaultErrKey, err)
	}
	writeChans := make([]chan model.Error, len(outputTargets))
	for i := 0; i < len(outputTargets); i++ {
		ch := make(chan model.Error)
		go p.dispatchWrite(
			ch,
			resourceData, resourceResult,
			formatters[i], writers[i],
			outputTargets[i],
		)
		writeChans[i] = ch
	}
	outputError := make(model.Error)
	for i := 0; i < len(outputTargets); i++ {
		err := <-writeChans[i]
		if err != nil {
			key := fmt.Sprintf("%s [%s]", defaultErrKey, outputTargets[i].Name)
			outputError[key] = err
		}
	}
	if len(outputError) > 0 {
		return outputError
	}
	return nil
}

func (p *Pipeline) dispatchWrite(
	ch chan model.Error,
	dataList []*model.Data, resultsList []string,
	formatter model.Format, writter model.Writer,
	t *model.OutputTarget) {
	const key = "dispatchWrite"
	writeErr := make(model.Error)
	for i := 0; i < len(dataList); i++ {
		content, err := formatter(dataList[i].Content)
		if err != nil {
			k := fmt.Sprintf("%s [%s]", key, dataList[i].Path)
			writeErr[k] = err
			continue
		}
		newData := &model.Data{
			Type:    t.Type,
			Path:    path.Join(t.Path, dataList[i].Path),
			Content: content,
		}
		if err := writter.Write(newData); err != nil {
			k := fmt.Sprintf("%s [%s]", key, dataList[i].Path)
			writeErr[k] = err
		}
	}
	if len(writeErr) > 0 {
		ch <- writeErr
	} else {
		ch <- nil
	}
}

func (p *Pipeline) getOutputWriter(target []*model.OutputTarget) ([]model.Writer, model.Error) {
	const defaultErrKey = "getOutputFormatter"
	outputWriter := make([]model.Writer, len(target))
	outputError := make(model.Error)
	for i, t := range target {
		fn, err := io.Writers.Get(t.Type)
		if err != nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			outputError[key] = err
			continue
		}
		outputWriter[i] = fn
	}
	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputWriter, nil
}

func (p *Pipeline) getOutputFormatter(target []*model.OutputTarget) ([]model.Format, model.Error) {
	const defaultErrKey = "getOutputFormatter"
	outputFormat := make([]model.Format, len(target))
	outputError := make(model.Error)
	for i, t := range target {
		fn, err := formatter.Formats.Get(jsonFormat, t.Format)
		if err != nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			outputError[key] = err
			continue
		}
		outputFormat[i] = fn
	}
	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputFormat, nil
}

func (p *Pipeline) dispatchEvaluate(evaluator *Evaluator, listOfData []*model.Data) []chan *evaluateWrapper {
	progress := NewProgress("dispatch evaluation", len(listOfData))
	evalChans := make([]chan *evaluateWrapper, len(listOfData))
	for i, data := range listOfData {
		ch := make(chan *evaluateWrapper)
		go func(c chan *evaluateWrapper, e *Evaluator, d *model.Data) {
			result, err := e.Evaluate(d)
			if err != nil {
				c <- &evaluateWrapper{
					Error: err,
				}
			} else {
				c <- &evaluateWrapper{
					Result: result,
				}
			}
		}(ch, evaluator, data)
		evalChans[i] = ch
		progress.Increment()
	}
	progress.Wait()
	return evalChans
}

func (p *Pipeline) dispatchValidate(validator *Validator, listOfData []*model.Data) []chan model.Error {
	progress := NewProgress("dispatch validation", len(listOfData))
	validateChans := make([]chan model.Error, len(listOfData))
	for i, data := range listOfData {
		ch := make(chan model.Error)
		go func(c chan model.Error, v *Validator, d *model.Data) {
			c <- v.Validate(d)
		}(ch, validator, data)
		validateChans[i] = ch
		progress.Increment()
	}
	progress.Wait()
	return validateChans
}
