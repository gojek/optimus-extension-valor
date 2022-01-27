package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gojek/optimus-extension-valor/model"
	_ "github.com/gojek/optimus-extension-valor/plugin/io" // init error writer
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const errorWriterType = "std"

var errorWriter model.Writer

func init() {
	writerFn, err := io.Writers.Get(errorWriterType)
	if err != nil {
		panic(err)
	}
	errorWriter = writerFn(model.TreatmentError)
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
	evaluate model.Evaluate,
	batchSize int,
	newProgress model.NewProgress,
) (*Pipeline, error) {
	if rcp == nil {
		return nil, errors.New("recipe is nil")
	}
	if evaluate == nil {
		return nil, errors.New("evaluate function is nil")
	}
	if batchSize < 0 {
		return nil, errors.New("batch size should be at least zero")
	}
	if newProgress == nil {
		return nil, errors.New("new progress function is nil")
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
func (p *Pipeline) Execute() error {
	for _, resourceRcp := range p.recipe.Resources {
		if err := p.validateFrameworkNames(resourceRcp); err != nil {
			return err
		}
		nameToFramework, err := p.getNameToFramework(resourceRcp)
		if err != nil {
			return err
		}
		nameToValidator, err := p.getNameToValidator(nameToFramework)
		if err != nil {
			return err
		}
		nameToEvaluator, err := p.getNameToEvaluator(nameToFramework)
		if err != nil {
			return err
		}
		if err := p.executeResource(resourceRcp, nameToValidator, nameToEvaluator); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) executeResource(resourceRcp *recipe.Resource, nameToValidator map[string]*Validator, nameToEvaluator map[string]*Evaluator) error {
	if resourceRcp == nil {
		return errors.New("resource recipe is nil")
	}
	resourcePaths, err := ExplorePaths(resourceRcp.Path, resourceRcp.Type, resourceRcp.Format)
	if err != nil {
		return err
	}
	outputError := make(model.Error)

	progress := p.newProgress(resourceRcp.Name, len(resourcePaths))
	batch := p.batchSize
	if batch == 0 || batch >= len(resourcePaths) {
		batch = len(resourcePaths)
	}
	counter := 0
	for counter < len(resourcePaths) {
		wg := &sync.WaitGroup{}
		mtx := &sync.Mutex{}

		for i := 0; i < batch && counter+i < len(resourcePaths); i++ {
			wg.Add(1)

			idx := counter + i
			go func(pt string, w *sync.WaitGroup, m *sync.Mutex) {
				defer w.Done()
				data, err := p.loader.LoadData(pt, resourceRcp.Type, resourceRcp.Format)
				if err != nil {
					m.Lock()
					outputError[pt] = err
					m.Unlock()
					return
				}
				for _, frameworkName := range resourceRcp.FrameworkNames {
					validator := nameToValidator[frameworkName]
					handleErr := func(success bool, err error) bool {
						if err != nil {
							var message string
							if e, ok := err.(model.Error); ok {
								message = string(e.JSON())
							} else {
								message = err.Error()
							}
							errorWriter.Write(&model.Data{
								Type:    "std",
								Path:    pt,
								Content: []byte(message),
							})
							return false
						}
						if !success {
							m.Lock()
							outputError[pt] = errors.New("business error encountered")
							m.Unlock()
							return false
						}
						return true
					}
					if validator != nil {
						success, err := validator.Validate(data)
						if ok := handleErr(success, err); !ok {
							return
						}
					}
					evaluator := nameToEvaluator[frameworkName]
					if evaluator != nil {
						success, err := evaluator.Evaluate(data)
						if ok := handleErr(success, err); !ok {
							return
						}
					}
				}
			}(resourcePaths[idx], wg, mtx)
		}
		wg.Wait()

		increment := batch
		if increment+counter > len(resourcePaths) {
			increment = len(resourcePaths) - counter
		}
		progress.Increase(increment)
		counter += batch
	}
	progress.Wait()
	if len(outputError) > 0 {
		return outputError
	}
	return nil
}

func (p *Pipeline) getNameToValidator(nameToFramework map[string]*model.Framework) (map[string]*Validator, error) {
	outputValidator := make(map[string]*Validator)
	outputError := make(model.Error)
	for name, framework := range nameToFramework {
		validator, err := NewValidator(framework)
		if err != nil {
			outputError[name] = err
		} else {
			outputValidator[name] = validator
		}
	}
	if len(outputError) > 0 {
		return nil, outputError
	}
	return outputValidator, nil
}

func (p *Pipeline) getNameToEvaluator(nameToFramework map[string]*model.Framework) (map[string]*Evaluator, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputEvaluator := make(map[string]*Evaluator)
	outputError := make(model.Error)
	for name, framework := range nameToFramework {
		wg.Add(1)

		go func(n string, f *model.Framework, w *sync.WaitGroup, m *sync.Mutex) {
			defer w.Done()
			evaluator, err := NewEvaluator(f, p.evaluate)
			if err != nil {
				m.Lock()
				outputError[n] = err
				m.Unlock()
			} else {
				m.Lock()
				outputEvaluator[n] = evaluator
				m.Unlock()
			}
		}(name, framework, wg, mtx)
	}
	if len(outputError) > 0 {
		return nil, outputError
	}
	wg.Wait()
	return outputEvaluator, nil
}

func (p *Pipeline) getNameToFramework(rcp *recipe.Resource) (map[string]*model.Framework, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	nameToFramework := make(map[string]*model.Framework)
	outputError := make(model.Error)
	for _, frameworkName := range rcp.FrameworkNames {
		wg.Add(1)

		go func(frameworkRcp *recipe.Framework, w *sync.WaitGroup, m *sync.Mutex) {
			defer w.Done()
			framework, err := p.loader.LoadFramework(frameworkRcp)
			if err != nil {
				m.Lock()
				outputError[frameworkRcp.Name] = err
				m.Unlock()
			} else {
				m.Lock()
				nameToFramework[frameworkRcp.Name] = framework
				m.Unlock()
			}
		}(p.nameToFrameworkRecipe[frameworkName], wg, mtx)
	}
	wg.Wait()

	if len(outputError) > 0 {
		return nil, outputError
	}
	return nameToFramework, nil
}

func (p *Pipeline) validateFrameworkNames(resourceRcp *recipe.Resource) error {
	outputError := make(model.Error)
	for _, frameworkName := range resourceRcp.FrameworkNames {
		if p.nameToFrameworkRecipe[frameworkName] == nil {
			outputError[frameworkName] = fmt.Errorf("not found for resource [%s]", resourceRcp.Name)
		}
	}
	if len(outputError) > 0 {
		return outputError
	}
	return nil
}
