package core

import (
	"errors"
	"fmt"
	"strings"
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
	newProgress model.NewProgress

	nameToFrameworkRecipe map[string]*recipe.Framework
}

// NewPipeline initializes pipeline process
func NewPipeline(
	rcp *recipe.Recipe,
	evaluate model.Evaluate,
	newProgress model.NewProgress,
) (*Pipeline, error) {
	if rcp == nil {
		return nil, errors.New("recipe is nil")
	}
	if evaluate == nil {
		return nil, errors.New("evaluate function is nil")
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
		newProgress:           newProgress,
		nameToFrameworkRecipe: nameToFrameworkRecipe,
	}, nil
}

// Execute executes pipeline process
func (p *Pipeline) Execute() error {
	for _, resourceRcp := range p.recipe.Resources {
		fmt.Printf("Resource [%s]\n", strings.ToUpper(resourceRcp.Name))
		fmt.Println("o> validating framework names")
		if err := p.validateFrameworkNames(resourceRcp); err != nil {
			return err
		}
		fmt.Println("o> loading the required framework data")
		nameToFramework, err := p.getFrameworkNameToFramework(resourceRcp)
		if err != nil {
			return err
		}
		fmt.Println("o> loading the required validator")
		nameToValidator, err := p.getFrameworkNameToValidator(nameToFramework)
		if err != nil {
			return err
		}
		fmt.Println("o> loading the required evaluator")
		nameToEvaluator, err := p.getFrameworkNameToEvaluator(nameToFramework)
		if err != nil {
			return err
		}
		fmt.Println("o> executing resource")
		if err := p.executeOnResource(resourceRcp, nameToValidator, nameToEvaluator); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func (p *Pipeline) executeOnResource(resourceRcp *recipe.Resource, nameToValidator map[string]*Validator, nameToEvaluator map[string]*Evaluator) error {
	if resourceRcp == nil {
		return errors.New("resource recipe is nil")
	}
	resourcePaths, err := ExplorePaths(resourceRcp.Path, resourceRcp.Type, resourceRcp.Format, resourceRcp.RegexPattern)
	if err != nil {
		return err
	}

	outputError := &model.Error{}
	handleErr := func(resourcePath, processType, frameworkName string, success bool, err error) bool {
		if err != nil {
			var message string
			if e, ok := err.(*model.Error); ok {
				message = string(e.JSON())
			} else {
				message = err.Error()
			}
			errorWriter.Write(&model.Data{
				Type:    "std",
				Path:    resourcePath,
				Content: []byte(message),
			})
			outputError.Add(resourcePath,
				fmt.Errorf("%s on framework [%s] encountered execution error",
					processType, frameworkName,
				),
			)
			return false
		}
		if !success {
			outputError.Add(resourcePath,
				fmt.Errorf("%s on framework [%s] encountered business error",
					processType, frameworkName,
				),
			)
			return false
		}
		return true
	}

	batch := resourceRcp.BatchSize
	if batch <= 0 || batch > len(resourcePaths) {
		batch = len(resourcePaths)
	}

	fmt.Printf(" [batch size: %d]\n", batch)
	progress := p.newProgress(resourceRcp.Name, len(resourcePaths))

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
					outputError.Add(pt, err)
					return
				}
				for _, frameworkName := range resourceRcp.FrameworkNames {
					validator := nameToValidator[frameworkName]
					if validator != nil {
						success, err := validator.Validate(data)
						if ok := handleErr(pt, "validation", frameworkName, success, err); !ok {
							return
						}
					}
					evaluator := nameToEvaluator[frameworkName]
					if evaluator != nil {
						success, err := evaluator.Evaluate(data)
						if ok := handleErr(pt, "evaluation", frameworkName, success, err); !ok {
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

	if outputError.Length() > 0 {
		return outputError
	}
	return nil
}

func (p *Pipeline) getFrameworkNameToValidator(nameToFramework map[string]*model.Framework) (map[string]*Validator, error) {
	outputValidator := make(map[string]*Validator)
	outputError := &model.Error{}
	for name, framework := range nameToFramework {
		validator, err := NewValidator(framework)
		if err != nil {
			outputError.Add(name, err)
		} else {
			outputValidator[name] = validator
		}
	}
	if outputError.Length() > 0 {
		return nil, outputError
	}
	return outputValidator, nil
}

func (p *Pipeline) getFrameworkNameToEvaluator(nameToFramework map[string]*model.Framework) (map[string]*Evaluator, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	outputEvaluator := make(map[string]*Evaluator)
	outputError := &model.Error{}
	for name, framework := range nameToFramework {
		wg.Add(1)

		go func(n string, f *model.Framework, w *sync.WaitGroup, m *sync.Mutex) {
			defer w.Done()
			evaluator, err := NewEvaluator(f, p.evaluate)
			if err != nil {
				outputError.Add(n, err)
			} else {
				m.Lock()
				outputEvaluator[n] = evaluator
				m.Unlock()
			}
		}(name, framework, wg, mtx)
	}
	wg.Wait()

	if outputError.Length() > 0 {
		return nil, outputError
	}
	return outputEvaluator, nil
}

func (p *Pipeline) getFrameworkNameToFramework(rcp *recipe.Resource) (map[string]*model.Framework, error) {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	nameToFramework := make(map[string]*model.Framework)
	outputError := &model.Error{}
	for _, frameworkName := range rcp.FrameworkNames {
		wg.Add(1)

		go func(frameworkRcp *recipe.Framework, w *sync.WaitGroup, m *sync.Mutex) {
			defer w.Done()
			framework, err := p.loader.LoadFramework(frameworkRcp)
			if err != nil {
				outputError.Add(frameworkRcp.Name, err)
			} else {
				m.Lock()
				nameToFramework[frameworkRcp.Name] = framework
				m.Unlock()
			}
		}(p.nameToFrameworkRecipe[frameworkName], wg, mtx)
	}
	wg.Wait()

	if outputError.Length() > 0 {
		return nil, outputError
	}
	return nameToFramework, nil
}

func (p *Pipeline) validateFrameworkNames(resourceRcp *recipe.Resource) error {
	outputError := &model.Error{}
	for _, frameworkName := range resourceRcp.FrameworkNames {
		if p.nameToFrameworkRecipe[frameworkName] == nil {
			outputError.Add(frameworkName, fmt.Errorf("not found for resource [%s]", resourceRcp.Name))
		}
	}
	if outputError.Length() > 0 {
		return outputError
	}
	return nil
}
