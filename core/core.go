package core

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/gojek/optimus-extension-valor/model"
	_ "github.com/gojek/optimus-extension-valor/plugin/io" // init error writer
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const (
	errorWriterType     = "std"
	defaultProgressType = "verbose"
)

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
	recipe      *recipe.Recipe
	loader      *Loader
	evaluate    model.Evaluate
	batchSize   int
	newProgress model.NewProgress

	nameToFrameworkRecipe map[string]*recipe.Framework
}

// NewPipeline initializes pipeline process
func NewPipeline(
	rcp *recipe.Recipe,
	batchSize int,
	evaluate model.Evaluate,
	newProgress model.NewProgress,
) (*Pipeline, model.Error) {
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
		evaluate:              evaluate,
		batchSize:             batchSize,
		newProgress:           newProgress,
		nameToFrameworkRecipe: nameToFrameworkRecipe,
	}, nil
}

// Execute executes pipeline process
func (p *Pipeline) Execute() model.Error {
	const defaultErrKey = "Execute"
	for _, resourceRcp := range p.recipe.Resources {
		decorate := strings.Repeat("=", 12)
		fmt.Printf("%s PROCESSING RESOURCE [%s] %s\n", decorate, resourceRcp.Name, decorate)

		fmt.Println("> Loading resource")
		resource, err := p.loader.LoadResource(resourceRcp)
		if err != nil {
			fmt.Println("* Loading failed!!!")
			return err
		}
		fmt.Printf("> Loading finished\n")

		for _, frameworkName := range resource.FrameworkNames {
			decorate := strings.Repeat(":", 5)
			fmt.Printf("%s Processing Framework [%s] %s\n", decorate, frameworkName, decorate)
			frameworkRcp := p.nameToFrameworkRecipe[frameworkName]

			fmt.Println(" >> Loading framework")
			framework, err := p.loader.LoadFramework(frameworkRcp)
			if err != nil {
				fmt.Println(" ** Loading failed!!!")
				return err
			}
			fmt.Printf(" >  Loading finished\n")

			fmt.Println(" >> Validating resource")
			success := p.executeValidate(framework, resource.ListOfData)
			if !success {
				fmt.Println(" ** Validation failed!!!")
				key := fmt.Sprintf("%s [validate: %s]", defaultErrKey, frameworkName)
				return model.BuildError(key, errors.New("error is met during validation"))
			}
			fmt.Printf(" >  Validation finished\n")

			fmt.Println(" >> Evaluating resource")
			success = p.executeEvaluate(framework, resource.ListOfData, framework.OutputTargets)
			if !success {
				fmt.Println(" ** Evaluation failed!!!")
				key := fmt.Sprintf("%s [evaluate: %s]", defaultErrKey, frameworkName)
				return model.BuildError(key, errors.New("error is met during evaluation"))
			}
			fmt.Printf(" >  Evaluation finished\n")
		}
		fmt.Println()
	}
	return nil
}

func (p *Pipeline) executeValidate(framework *model.Framework, resourceData []*model.Data) bool {
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
	progress := p.newProgress(fmt.Sprintf("%s [%s]", defaultErrKey, framework.Name), len(resourceData))
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
			if success {
				m.Lock()
				progress.Increment()
				m.Unlock()
			}
		}(i, validator, wg, mtx, data)
	}
	wg.Wait()
	progress.Wait()
	return success
}

func (p *Pipeline) executeEvaluate(framework *model.Framework, resourceData []*model.Data, outputTargets []*model.OutputTarget) bool {
	const defaultErrKey = "evaluate"
	evaluator, err := NewEvaluator(framework, p.evaluate)
	if err != nil {
		errorWriter.Write(&model.Data{
			Type:    errorWriterType,
			Content: err.JSON(),
			Path:    defaultErrKey,
		})
		return false
	}
	progress := p.newProgress(fmt.Sprintf("%s [%s]", defaultErrKey, framework.Name), len(resourceData))

	batch := p.batchSize
	if batch <= 0 || batch >= len(resourceData) {
		batch = len(resourceData)
	}
	counter := 0

	success := true
	for counter < len(resourceData) {
		wg := &sync.WaitGroup{}
		mtx := &sync.Mutex{}
		for i := 0; i < batch && counter+i < len(resourceData); i++ {
			wg.Add(1)
			data := resourceData[counter+i]
			go func(idx int, w *sync.WaitGroup, m *sync.Mutex, d *model.Data) {
				defer w.Done()
				rst, err := evaluator.Evaluate(d)
				currentSuccess := true
				if err != nil {
					pt := fmt.Sprintf("%s [%d]", defaultErrKey, idx)
					if d != nil {
						pt = d.Path
					}
					errorWriter.Write(&model.Data{
						Type:    errorWriterType,
						Content: err.JSON(),
						Path:    pt,
					})
					currentSuccess = false
				} else {
					result := &model.Data{
						Type:    d.Type,
						Path:    d.Path,
						Content: []byte(rst),
					}
					currentSuccess = p.writeOutput(result, outputTargets)
				}
				if !currentSuccess {
					m.Lock()
					success = false
					m.Unlock()
				}
			}(i, wg, mtx, data)
		}
		wg.Wait()
		for i := 0; i < batch && counter+i < len(resourceData); i++ {
			progress.Increment()
		}
		counter += batch
	}
	progress.Wait()
	return success
}

func (p *Pipeline) writeOutput(result *model.Data, outputTargets []*model.OutputTarget) bool {
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
			newContent, err := formatters[idx](result.Content)
			if err != nil {
				errorWriter.Write(&model.Data{
					Type:    result.Type,
					Path:    result.Path,
					Content: err.JSON(),
				})
				m.Lock()
				success = false
				m.Unlock()
			}
			if !success {
				return
			}
			newResult := &model.Data{
				Type:    result.Type,
				Path:    path.Join(outputTargets[idx].Path, result.Path),
				Content: newContent,
			}
			if err := writers[idx].Write(newResult); err != nil {
				m.Lock()
				success = false
				m.Unlock()
				errorWriter.Write(&model.Data{
					Type:    errorWriterType,
					Path:    fmt.Sprintf("%s [%d]", defaultErrKey, idx),
					Content: err.JSON(),
				})
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
