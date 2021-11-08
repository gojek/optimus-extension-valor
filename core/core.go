package core

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gojek/optimus-extension-valor/model"
	_ "github.com/gojek/optimus-extension-valor/plugin/io" // init error writer
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"

	"github.com/google/go-jsonnet"
)

const errorWriterType = "std"

var errorWriter model.Writer

func init() {
	writer, err := io.Writers.Get(errorWriterType)
	if err != nil {
		panic(err)
	}
	errorWriter = writer
}

const (
	jsonFormat    = "json"
	jsonnetFormat = "jsonnet"
)

var skipReformat = map[string]bool{
	jsonFormat:    true,
	jsonnetFormat: true,
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

		start := time.Now()
		fmt.Println("> Loading resource")
		resource, err := p.loader.LoadResource(resourceRcp)
		if err != nil {
			fmt.Println("* Loading failed!!!")
			return err
		}
		fmt.Printf("> Loading finished [%v]\n", time.Now().Sub(start))

		for _, frameworkName := range resource.FrameworkNames {
			decorate := strings.Repeat(":", 5)
			fmt.Printf("%s Processing Framework [%s] %s\n", decorate, frameworkName, decorate)
			frameworkRcp := p.nameToFrameworkRecipe[frameworkName]

			start := time.Now()
			fmt.Println(" >> Loading framework")
			framework, err := p.loader.LoadFramework(frameworkRcp)
			if err != nil {
				fmt.Println(" ** Loading failed!!!")
				return err
			}
			fmt.Printf(" >  Loading finished [%v]\n", time.Now().Sub(start))

			start = time.Now()
			fmt.Println(" >> Validating resource")
			success := p.validate(framework, resource.ListOfData)
			if !success {
				fmt.Println(" ** Validation failed!!!")
				key := fmt.Sprintf("%s [validate: %s]", defaultErrKey, frameworkName)
				return model.BuildError(key, errors.New("error is met during validation"))
			}
			fmt.Printf(" >  Validation finished [%v]\n", time.Now().Sub(start))

			start = time.Now()
			fmt.Println(" >> Evaluating resource")
			evalResults, success := p.evaluate(p.vm, framework, resource.ListOfData)
			if !success {
				fmt.Println(" ** Evaluation failed!!!")
				key := fmt.Sprintf("%s [evaluate: %s]", defaultErrKey, frameworkName)
				return model.BuildError(key, errors.New("error is met during evaluation"))
			}
			fmt.Printf(" >  Evaluation finished [%v]\n", time.Now().Sub(start))

			start = time.Now()
			fmt.Println(" >> Writing result")
			success = p.writeOutput(evalResults, framework.OutputTargets)
			if !success {
				fmt.Println(" ** Writing failed!!!")
				key := fmt.Sprintf("%s [write: %s]", defaultErrKey, frameworkName)
				return model.BuildError(key, errors.New("error is met during write"))
			}
			fmt.Printf(" >  Writing finished [%v]\n", time.Now().Sub(start))
		}
		fmt.Println()
	}
	return nil
}

func (p *Pipeline) validate(framework *model.Framework, resourceData []*model.Data) bool {
	const defaultErrKey = "validate"
	validator, err := NewValidator(framework)
	if err != nil {
		errorWriter.Write(&model.Data{
			Type:    errorWriterType,
			Content: err.JSON(),
			Path:    defaultErrKey,
		})
		return false
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	success := true
	for i, data := range resourceData {
		wg.Add(1)

		go func(idx int, v *Validator, w *sync.WaitGroup, m *sync.Mutex, d *model.Data) {
			defer w.Done()

			if err := v.Validate(d); err != nil {
				m.Lock()
				success = false
				m.Unlock()

				pt := fmt.Sprintf("%s [%d]", defaultErrKey, idx)
				if d != nil {
					pt = fmt.Sprintf("%s [%s]", defaultErrKey, d.Path)
				}
				errorWriter.Write(&model.Data{
					Type:    errorWriterType,
					Content: err.JSON(),
					Path:    pt,
				})
			}
		}(i, validator, wg, mtx, data)
	}
	wg.Wait()
	return success
}

func (p *Pipeline) evaluate(vm *jsonnet.VM, framework *model.Framework, resourceData []*model.Data) ([]*model.Data, bool) {
	const defaultErrKey = "evaluate"
	evaluator, err := NewEvaluator(framework, vm)
	if err != nil {
		errorWriter.Write(&model.Data{
			Type:    errorWriterType,
			Content: err.JSON(),
			Path:    defaultErrKey,
		})
		return nil, false
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputResult := make([]*model.Data, len(resourceData))
	success := true
	for i, data := range resourceData {
		wg.Add(1)

		go func(idx int, v *Evaluator, w *sync.WaitGroup, m *sync.Mutex, d *model.Data) {
			defer w.Done()

			rst, err := v.Evaluate(d)
			if err != nil {
				m.Lock()
				success = false
				m.Unlock()

				pt := fmt.Sprintf("%s [%d]", defaultErrKey, idx)
				if d != nil {
					pt = d.Path
				}
				errorWriter.Write(&model.Data{
					Type:    errorWriterType,
					Content: err.JSON(),
					Path:    pt,
				})
			} else {
				if success {
					m.Lock()
					outputResult[idx] = &model.Data{
						Type:    d.Type,
						Path:    d.Path,
						Content: []byte(rst),
					}
					m.Unlock()
				}
			}
		}(i, evaluator, wg, mtx, data)
	}
	wg.Wait()

	if success {
		return outputResult, true
	}
	return nil, false
}

func (p *Pipeline) writeOutput(evalResults []*model.Data, outputTargets []*model.OutputTarget) bool {
	const defaultErrKey = "writeOutput"
	formatters, err := p.getOutputFormatter(outputTargets)
	if err != nil {
		errorWriter.Write(&model.Data{
			Type:    errorWriterType,
			Content: err.JSON(),
			Path:    defaultErrKey,
		})
		return false
	}
	writers, err := p.getOutputWriter(outputTargets)
	if err != nil {
		errorWriter.Write(&model.Data{
			Type:    errorWriterType,
			Content: err.JSON(),
			Path:    defaultErrKey,
		})
		return false
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	success := true
	for i := 0; i < len(outputTargets); i++ {
		wg.Add(1)

		go func(idx int, w *sync.WaitGroup, m *sync.Mutex) {
			defer w.Done()

			newResults := make([]*model.Data, len(evalResults))
			currentSuccess := true
			for j := 0; j < len(evalResults); j++ {
				newContent, err := formatters[idx](evalResults[j].Content)
				if err != nil {
					errorWriter.Write(&model.Data{
						Type:    evalResults[j].Path,
						Path:    evalResults[j].Path,
						Content: err.JSON(),
					})
					currentSuccess = false
					continue
				}
				newResults[j] = &model.Data{
					Type:    evalResults[j].Type,
					Path:    evalResults[j].Path,
					Content: newContent,
				}
			}
			if currentSuccess {
				if err := writers[idx].Write(newResults...); err != nil {
					currentSuccess = false
					errorWriter.Write(&model.Data{
						Type:    errorWriterType,
						Path:    fmt.Sprintf("%s [%d]", defaultErrKey, idx),
						Content: err.JSON(),
					})
				}
			}
			if !currentSuccess {
				m.Lock()
				success = false
				m.Unlock()
			}
		}(i, wg, mtx)
	}
	wg.Wait()
	return success
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
