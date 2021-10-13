package core

import (
	"errors"
	"fmt"
	"strings"
)

// Snippet contains snippet for all Procedures in a Framework
type Snippet struct {
	nameToIsRegistered  map[string]bool
	nameToProcedure     map[string]string
	nameToOutputIsError map[string]bool
}

// GetByProcedure gets a snippet based on a Procedure name
func (s *Snippet) GetByProcedure(name string) (string, error) {
	if !s.nameToIsRegistered[name] {
		return "", fmt.Errorf("procedure [%s] is not registered", name)
	}
	return s.nameToProcedure[name], nil
}

// OutputIsError gets whether Output from a Procedure, specified by its name, is
// considered an error or not
func (s *Snippet) OutputIsError(name string) (bool, error) {
	if !s.nameToIsRegistered[name] {
		return false, fmt.Errorf("procedure [%s] is not registered", name)
	}
	return s.nameToOutputIsError[name], nil
}

// ListProcedureNames list all registered Procedure names
func (s *Snippet) ListProcedureNames() []string {
	var output []string
	for name := range s.nameToIsRegistered {
		output = append(output, name)
	}
	return output
}

// NewSnippet initializes a new Snippet for a Framework
func NewSnippet(framework *Framework) (*Snippet, error) {
	if framework == nil {
		return nil, errors.New("framework is nil")
	}
	definitionSnippet := buildDefinitionSnippet(framework.definitions)

	nameToIsRegistered := make(map[string]bool)
	nameToProcedure := make(map[string]string)
	nameToOutputIsError := make(map[string]bool)
	for _, procedure := range framework.procedures {
		temporarySnippet := ""
		for _, data := range procedure.data {
			temporarySnippet += string(data.Content) + "\n"
		}

		procedureSnippet := definitionSnippet + "\n" + temporarySnippet
		procedureName := procedure.name
		outputIsError := procedure.outputIsError

		nameToProcedure[procedureName] = procedureSnippet
		nameToOutputIsError[procedureName] = outputIsError
		nameToIsRegistered[procedureName] = true
	}
	return &Snippet{
		nameToIsRegistered:  nameToIsRegistered,
		nameToProcedure:     nameToProcedure,
		nameToOutputIsError: nameToOutputIsError,
	}, nil
}

func buildDefinitionSnippet(definitions []*Container) string {
	var output string
	for _, def := range definitions {
		varName := def.name
		var values []string
		for _, d := range def.data {
			values = append(values, string(d.Content))
		}
		varValues := fmt.Sprintf("[%s]", strings.Join(values, ","))
		output += fmt.Sprintf("local %s = %s;\n", varName, varValues)
	}
	return output
}
