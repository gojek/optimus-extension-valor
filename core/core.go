package core

import (
	"errors"
	"fmt"
	"log"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/recipe"
	"github.com/gojek/optimus-extension-valor/registry/formatter"
	"github.com/gojek/optimus-extension-valor/registry/io"
)

const (
	defaultFormat = "json"
	skipResult    = "null\n"
)

var isNotToFormat = map[string]bool{
	"json":    true,
	"jsonnet": true,
}

// Pipeline holds information on how a pipeline is executed
type Pipeline struct {
	recipe *recipe.Recipe

	loaded    bool
	built     bool
	evaluated bool

	nameToFramework        map[string]*Framework
	frameworkNameToSnippet map[string]*Snippet

	resourceList []*Resource

	schemaRstList       []*schemaResult
	procedureRstList    []*procedureResult
	pipelineBusinessErr bool
	pipelineExecErr     bool
}

// Flush flushes pipeline results based on output
func (p *Pipeline) Flush() error {
	if !p.evaluated {
		return errors.New("pipeline is not executed")
	}

	log.Println("Flushing Output")
	if len(p.schemaRstList) > 0 {
		schemaOutputMap := p.schemaRstListToOutputMap(p.schemaRstList)
		err := p.flushOutput("schema", schemaOutputMap)
		if err != nil {
			return err
		}
	}
	procedureOutputMap := p.procedureRstListToOutputMap(p.procedureRstList)
	if len(p.procedureRstList) > 0 {
		err := p.flushOutput("procedure", procedureOutputMap)
		if err != nil {
			return err
		}
	}

	if p.pipelineExecErr {
		log.Println("Finished with Execution Error")
		return errors.New("execution error encountered")
	} else if p.pipelineBusinessErr {
		log.Println("Finished with Business Error")
		return errors.New("business error encountered")
	} else {
		log.Println("Finished Successfully")
	}
	fmt.Println()
	return nil
}

func (p *Pipeline) flushOutput(key string, outputMap map[*recipe.Metadata][]*model.Data) error {
	var outputErrors []error
	progress := NewProgress(key, len(outputMap))
	for output, listOfData := range outputMap {
		updatedListOfData := make([]*model.Data, len(listOfData))
		for i, d := range listOfData {
			formatterFn, err := formatter.Formatters.Get(defaultFormat, output.Format)
			if err != nil {
				newErr := fmt.Errorf("getting formatter: %v", err)
				outputErrors = append(outputErrors, newErr)
				continue
			}
			content, err := formatterFn(d.Content)
			if err != nil {
				newErr := fmt.Errorf("formatting content from [%s] to [%s]: %v",
					defaultFormat, output.Format, err,
				)
				outputErrors = append(outputErrors, newErr)
				continue
			}
			updatedListOfData[i] = &model.Data{
				Path:     d.Path,
				Content:  content,
				Metadata: d.Metadata,
			}
		}
		if len(outputErrors) > 0 {
			continue
		}
		writerFn, err := io.Writers.Get(output.Type)
		if err != nil {
			newErr := fmt.Errorf("getting writer [%s]: %v", output.Type, err)
			outputErrors = append(outputErrors, newErr)
			continue
		}
		writer := writerFn(output.Path, map[string]string{
			"type":      output.Type,
			"name":      output.Name,
			"path":      output.Path,
			"extension": output.Format,
		})
		err = writer.Write(updatedListOfData...)
		if err != nil {
			newErr := fmt.Errorf("writing list of data: %v", err)
			outputErrors = append(outputErrors, newErr)
		}
		progress.Increment()
	}
	progress.Wait()
	if len(outputErrors) > 0 {
		printErrors(outputErrors)
		return errors.New("one or more errors are encountered")
	}
	return nil
}

func (p *Pipeline) procedureRstListToOutputMap(procedureRstList []*procedureResult) map[*recipe.Metadata][]*model.Data {
	outputMap := make(map[*recipe.Metadata][]*model.Data)
	for _, procedureRst := range procedureRstList {
		output := p.resultToOutput(procedureRst)
		for _, o := range procedureRst.outputTargets {
			outputMap[o] = append(outputMap[o], output)
		}
	}
	return outputMap
}

func (p *Pipeline) schemaRstListToOutputMap(schemaRstList []*schemaResult) map[*recipe.Metadata][]*model.Data {
	outputMap := make(map[*recipe.Metadata][]*model.Data)
	for _, schemaRst := range schemaRstList {
		output := p.resultToOutput(schemaRst)
		for _, o := range schemaRst.outputTargets {
			outputMap[o] = append(outputMap[o], output)
		}
	}
	return outputMap
}

func (p *Pipeline) resultToOutput(result interface{}) *model.Data {
	var path string
	var content []byte
	var metadata map[string]string
	if schemaRst, ok := result.(*schemaResult); ok {
		path = schemaRst.record.Path
		metadata = schemaRst.record.Metadata
		if schemaRst.err != nil {
			content = []byte(fmt.Sprintf("%#v", schemaRst.err))
		} else {
			content = schemaRst.result
		}
	} else {
		procedureRst := result.(*procedureResult)
		path = procedureRst.record.Path
		metadata = procedureRst.record.Metadata
		if procedureRst.err != nil {
			content = []byte(fmt.Sprintf("%#v", procedureRst.err))
		} else {
			content = procedureRst.result
		}
	}
	return &model.Data{
		Path:     path,
		Content:  content,
		Metadata: metadata,
	}
}

// Execute executes pipeline process
func (p *Pipeline) Execute() error {
	if !p.built {
		return errors.New("pipeline is not built")
	}

	log.Println("Evaluating Resource")
	p.evaluated = true
	var schemaRstList []*schemaResult
	var procedureRstList []*procedureResult
	pipelineExecErr := false
	pipelineBusinessErr := false
	for _, resource := range p.resourceList {
		progress := NewProgress(resource.container.name, len(resource.frameworkNames))
		resourceExecErr := false
		resourceBusinessErr := false
		for _, frwrkName := range resource.frameworkNames {
			framework := p.nameToFramework[frwrkName]
			if framework == nil {
				return fmt.Errorf("framework [%s] for resource [%s] is unknown", frwrkName, resource.container.name)
			}
			schemaExecErr := false
			schemaBusinessErr := false
			for _, schema := range framework.schemas {
				for _, schemaData := range schema.data {
					for _, record := range resource.container.data {
						rst, err := ValidateSchema(schemaData.Content, record.Content)
						schemaRst := &schemaResult{
							record:        record,
							schema:        schemaData,
							result:        rst,
							err:           err,
							outputTargets: framework.outputTargets,
						}
						if err != nil {
							schemaExecErr = true
						}
						if len(rst) > 0 {
							if schema.outputIsError {
								schemaBusinessErr = true
							}
							schemaRstList = append(schemaRstList, schemaRst)
						}
					}
				}
			}
			if schemaExecErr {
				resourceExecErr = true
				break
			}
			if schemaBusinessErr {
				resourceBusinessErr = true
				if !framework.allowError {
					break
				}
			}

			snippet := p.frameworkNameToSnippet[framework.name]
			procedureExecErr := false
			procedureBusinessErr := false
			for _, procedure := range framework.procedures {
				procedureSnippet, err := snippet.GetByProcedure(procedure.name)
				if err != nil {
					return err
				}
				for _, record := range resource.container.data {
					recordSnippet := fmt.Sprintf("local %s = %s;", resource.container.name, string(record.Content))
					completeSnippet := recordSnippet + "\n" + procedureSnippet
					rst, err := Evaluate(resource.container.name, completeSnippet)
					procedureRst := &procedureResult{
						record:        record,
						snippet:       completeSnippet,
						result:        rst,
						err:           err,
						outputTargets: framework.outputTargets,
					}
					if err != nil {
						procedureExecErr = true
					}
					if len(rst) > 0 && string(rst) != skipResult {
						outputIsError, err := snippet.OutputIsError(procedure.name)
						if err != nil {
							return err
						}
						if outputIsError {
							procedureBusinessErr = true
						}
						procedureRstList = append(procedureRstList, procedureRst)
					}
				}
			}
			if procedureExecErr {
				resourceExecErr = true
				break
			}
			if procedureBusinessErr {
				resourceBusinessErr = true
				if !framework.allowError {
					break
				}
			}
			progress.Increment()
		}
		progress.Wait()
		if resourceExecErr {
			pipelineExecErr = true
			break
		}
		if resourceBusinessErr {
			pipelineBusinessErr = true
			break
		}
	}

	p.schemaRstList = schemaRstList
	p.procedureRstList = procedureRstList
	p.pipelineExecErr = pipelineExecErr
	p.pipelineBusinessErr = pipelineBusinessErr

	log.Println("Finished evaluating")
	fmt.Println()
	return nil
}

// Build builds the pipeline
func (p *Pipeline) Build() error {
	if p.built {
		return errors.New("pipeline is already built")
	}

	log.Println("Building Pipeline")
	progress := NewProgress("building", len(p.nameToFramework))
	frameworkNameToSnippet := make(map[string]*Snippet)
	var buildingErrors []error
	for name, framework := range p.nameToFramework {
		snippet, err := NewSnippet(framework)
		if err != nil {
			buildingErrors = append(buildingErrors, err)
			continue
		}
		frameworkNameToSnippet[name] = snippet
		progress.Increment()
	}
	progress.Wait()

	if len(buildingErrors) > 0 {
		log.Println("Errors encountered")
		printErrors(buildingErrors)
		return errors.New("one or more errors are encountered")
	}
	fmt.Println()

	p.frameworkNameToSnippet = frameworkNameToSnippet
	p.built = true
	return nil
}

// Load loads the required data based on recipe
func (p *Pipeline) Load() error {
	if p.loaded {
		return errors.New("pipeline is already loaded")
	}

	log.Println("Loading Pipeline")
	resourceList, err := p.loadResourceList(p.recipe.Resources)
	if err != nil {
		return err
	}
	frameworkNameToRequired := p.getFrameworkNameToRequired(p.recipe.Resources)
	nameToFramework, err := p.loadNameToFramework(p.recipe.Frameworks, frameworkNameToRequired)
	if err != nil {
		return err
	}
	log.Println("Finished Loading")
	fmt.Println()

	p.resourceList = resourceList
	p.nameToFramework = nameToFramework
	p.loaded = true
	return nil
}

func (p *Pipeline) loadNameToFramework(rcpFrameworks []*recipe.Framework, nameToRequired map[string]bool) (map[string]*Framework, error) {
	progress := NewProgress("framework", len(rcpFrameworks))
	var frameworkNames []string
	nameToFramework := make(map[string]*Framework)
	var frameworkErrors []error
	for _, rcp := range rcpFrameworks {
		progress.Increment()
		if !nameToRequired[rcp.Name] {
			continue
		}
		frmwrk, err := LoadFramework(rcp)
		if err != nil {
			newErr := fmt.Errorf("framework [%s]: %v", rcp.Name, err)
			frameworkErrors = append(frameworkErrors, newErr)
			continue
		}
		frameworkNames = append(frameworkNames, frmwrk.name)
		nameToFramework[frmwrk.name] = frmwrk
	}
	progress.Wait()
	if len(frameworkErrors) > 0 {
		log.Println("Errors encountered")
		printErrors(frameworkErrors)
		return nil, errors.New("one or more errors are encountered")
	}
	return nameToFramework, nil
}

func (p *Pipeline) getFrameworkNameToRequired(rcpResources []*recipe.Resource) map[string]bool {
	output := make(map[string]bool)
	for _, c := range rcpResources {
		for _, name := range c.FrameworkNames {
			output[name] = true
		}
	}
	return output
}

func (p *Pipeline) loadResourceList(rcpResources []*recipe.Resource) ([]*Resource, error) {
	progress := NewProgress("resource", len(rcpResources))
	var resourceList []*Resource
	var resourceErrors []error
	for _, c := range rcpResources {
		resource, err := LoadResource(c)
		if err != nil {
			newErr := fmt.Errorf("resource [%s]: %v", c.Name, err)
			resourceErrors = append(resourceErrors, newErr)
			continue
		}
		resourceList = append(resourceList, resource)
		progress.Increment()
	}
	progress.Wait()
	if len(resourceErrors) > 0 {
		log.Println("Errors encountered")
		printErrors(resourceErrors)
		return nil, errors.New("one or more errors are encountered")
	}
	return resourceList, nil
}

// NewPipeline initializes pipeline process
func NewPipeline(rcp *recipe.Recipe) (*Pipeline, error) {
	if rcp == nil {
		return nil, errors.New("recipe is nil")
	}
	return &Pipeline{
		recipe: rcp,
	}, nil
}

func printErrors(errs []error) {
	for i, e := range errs {
		fmt.Printf("[%d] %v\n", i, e)
	}
}
