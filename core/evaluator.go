package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"
	"github.com/gojek/optimus-extension-valor/registry/endec"

	"github.com/google/go-jsonnet"
)

// Evaluator contains information on how to evaluate a Resource
type Evaluator struct {
	vm *jsonnet.VM

	framework         *model.Framework
	definitionSnippet string
}

// NewEvaluator initializes Evaluator
func NewEvaluator(framework *model.Framework, vm *jsonnet.VM) (*Evaluator, model.Error) {
	const defaultErrKey = "NewEvaluator"
	if framework == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("framework is nil"))
	}
	if vm == nil {
		return nil, model.BuildError(defaultErrKey, errors.New("vm is nil"))
	}
	definitionSnippet, err := buildAllDefinitions(vm, framework.Definitions)
	if err != nil {
		return nil, model.BuildError(defaultErrKey, err)
	}
	return &Evaluator{
		vm:                vm,
		framework:         framework,
		definitionSnippet: definitionSnippet,
	}, nil
}

// Evaluate evaluates snippet for a Resource data
func (e *Evaluator) Evaluate(resourceData *model.Data) (string, model.Error) {
	const defaultErrKey = "Evaluate"
	if resourceData == nil {
		return model.SkipNullValue, model.BuildError(defaultErrKey, errors.New("resource data is nil"))
	}
	resourceSnippet := string(resourceData.Content)
	previousOutputSnippet := model.SkipNullValue
	for _, procedure := range e.framework.Procedures {
		snippet, err := buildSnippet(resourceSnippet, e.definitionSnippet, previousOutputSnippet, procedure)
		if err != nil {
			return model.SkipNullValue, model.BuildError(defaultErrKey, err)
		}
		result, evalErr := e.vm.EvaluateAnonymousSnippet(procedure.Name, snippet)
		if evalErr != nil {
			return model.SkipNullValue, model.BuildError(defaultErrKey, evalErr)
		}
		if model.IsSkipResult[result] {
			previousOutputSnippet = model.SkipNullValue
		} else {
			evalResult := e.resultToError(result)
			if procedure.OutputIsError {
				return model.SkipNullValue, model.BuildError(defaultErrKey, evalResult)
			}
			previousOutputSnippet = result
		}
	}
	return previousOutputSnippet, nil
}

func (e *Evaluator) resultToError(result string) model.Error {
	const defaultErrKey = "resultToError"
	fn, err := endec.Decodes.Get(jsonFormat)
	if err != nil {
		return model.BuildError(defaultErrKey, err)
	}
	var output map[string]interface{}
	if err = fn([]byte(result), &output); err != nil {
		return model.BuildError(defaultErrKey, errors.New(result))
	}
	return e.castEvaluationResult(output)
}

func (e *Evaluator) castEvaluationResult(result map[string]interface{}) model.Error {
	output := make([]model.Error, len(result))
	for key, value := range result {
		if castedValue, ok := value.(map[string]interface{}); ok {
			childErr := e.castEvaluationResult(castedValue)
			output = append(output, model.BuildError(key, childErr))
		} else {
			output = append(output, model.BuildError(key, fmt.Errorf("%v", value)))
		}
	}
	return model.CombineErrors(output...)
}

func buildAllDefinitions(vm *jsonnet.VM, definitions []*model.Definition) (string, model.Error) {
	const defaultErrKey = "buildAllDefinitions"
	nameToSnippet := make(map[string]string)
	errOutput := make(model.Error)
	for i, def := range definitions {
		defSnippet, err := buildOneDefinition(vm, def)
		if err != nil {
			key := fmt.Sprintf("%s [%d]", defaultErrKey, i)
			if def != nil {
				key = fmt.Sprintf("%s [%s]", defaultErrKey, def.Name)
			}
			errOutput[key] = err
			continue
		}
		nameToSnippet[def.Name] = defSnippet
	}
	if len(errOutput) > 0 {
		return model.SkipNullValue, errOutput
	}
	var outputSnippets []string
	for key, value := range nameToSnippet {
		outputSnippets = append(outputSnippets, fmt.Sprintf(`"%s": %s`, key, value))
	}
	return fmt.Sprintf("{%s}", strings.Join(outputSnippets, "\n")), nil
}

func buildOneDefinition(vm *jsonnet.VM, definition *model.Definition) (string, model.Error) {
	const defaultErrKey = "buildOneDefinition"
	if definition == nil {
		return model.SkipNullValue, model.BuildError(defaultErrKey, errors.New("definition is nil"))
	}
	var defData string
	for i, data := range definition.ListOfData {
		if i > 0 {
			defData += ",\n"
		}
		defData += string(data.Content)
	}
	defSnippet := fmt.Sprintf("[%s]", defData)
	if definition.FunctionData != nil {
		defSnippet = fmt.Sprintf(`
/*
The line below is a place to put all data of a definition. The definitions
here is the one taken from "definition.Path" if it is specified.

The generated definition will follow:
	local definition = [...]; // an array

Detail:
	* variable name: "definition"
	* variable value: array of object

Example:
	local definition = [
		{ "name": "lion" },
		{ "name": "cow" },
	];
*/
local definition = %s;

/*
The line below is a place to put the function for definition. This function
is usually for constructing a new definition based on the definition
above. The function is provided by the user.

The format should follow:
	local construct(definition) = {...}; // an object
	or
	local construct(definition) = [...]; // an array

Detail:
	* function name: "construct"
	* argument: definition generated from the "definition.Path"
	* return: an object or an array

Example:
	local construct(definition) = [
		{ name: d.name, isChecked: true }
		for d in definition
	];
The example above adds field "isChecked" for each definition.
*/
%s

/*
The line below is to call the defined function. It requires the definition
being defined. The definition is then passed as argument. The user does
not need to specify it.

The format should follow:
	construct(definition)
*/
construct (definition)

`,
			defSnippet,
			string(definition.FunctionData.Content),
		)
		result, err := vm.EvaluateAnonymousSnippet(definition.Name, defSnippet)
		if err != nil {
			return model.SkipNullValue, model.BuildError(defaultErrKey, err)
		}
		defSnippet = result
	}
	return defSnippet, nil
}

func buildSnippet(
	resourceSnippet string,
	definitionSnippet string,
	previousOutputSnippet string,
	procedure *model.Procedure,
) (string, model.Error) {
	const defaultErrKey = "buildSnippet"
	if procedure == nil {
		return model.SkipNullValue, model.BuildError(defaultErrKey, errors.New("procedure is nil"))
	}
	output := fmt.Sprintf(`
/*
The line below is to declare the resource to be evaluated. This is auto-generated
taken from "resource.Path". The user only needs to define the resource itself.

The generated resource follow:
	local resource = {...}; // an object

Detail:
	* variable name: "resource"
	* variabl value: an object

Example:
	local resource = {"name": "unknown"};
*/
local resource = %s;

/*
The line below is to declare the definition. The definition is taken from the
defined definition, be it directly from the "definition.Path" or one that is
constructed. Since there can be multiple definitions being provided by the user,
with different name, then the definition here will be an object. This object
has key, which is taken from "definition.Name", and the value of object or array.

The format will follow:
	local definition = {
		"definition.Name": {...} // an object
		or
		"definition.Name": [...] // an array
	} // should be object

Detail:
	* variable name: "definition"
	* variable value: object of object or object of array

Example:
	local definition = {
		"animals": [...]
	}
*/
local definition = %s;

/*
The line below is to declare the previous output. Previous output
should be declared, even if the previous procedure does return
empty or even if the prevous result does not return anything.
If such case is encountered, then it should be set to be null value.

The format will follow:
	local previousOutput = {...}; // an object
	or
	local previousOutput = [...]; // an array

Detail:
	* variable name: "previousOutput"
	* varialbe value: an object or an array
*/
local previousOutput = %s;

/*
Line below is to declare the procedure defined by the user.

The format should follow:
	local evaluate(resource, definition, previousOutput) = {}; // should return object
	or
	local evaluate(resource, definition, previousOutput) = []; // should return array

Detail:
	* procedure name: "evaluate"
	* first argument: resource to be evaluated
	* second argument: definition to be used for evaluation
	* third argument: previous output
	* return: an objet or an array

Example:
	local evaluate(resource, definition, previousOutput) = {
		"valid": true,
	}
*/
%s

/*
The line below is to call the defined procedure. The user does not need to
define it as it is generated.

The format should follow:
	evaluate(resource, definition, previousOutput)
*/
evaluate (resource, definition, previousOutput)

`,
		resourceSnippet,
		definitionSnippet,
		previousOutputSnippet,
		string(procedure.Data.Content),
	)
	return output, nil
}
