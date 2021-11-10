package recipe

import (
	"fmt"
	"strings"

	"github.com/gojek/optimus-extension-valor/model"

	"github.com/go-playground/validator/v10"
)

const defaultValidateKey = "Validate"

// Validate validates the recipe
func Validate(rcp *Recipe) model.Error {
	if err := validator.New().Struct(rcp); err != nil {
		return model.BuildError(defaultValidateKey, err)
	}
	if err := validateAllResources(rcp.Resources); err != nil {
		return err
	}
	return validateAllFrameworks(rcp.Frameworks)
}

func validateAllResources(rcps []*Resource) model.Error {
	nameEncountered := make(map[string]int)
	for _, resourceRcp := range rcps {
		if err := ValidateResource(resourceRcp); err != nil {
			return err
		}
		nameEncountered[resourceRcp.Name]++
	}
	var duplicateNames []string
	for name, count := range nameEncountered {
		if count > 1 {
			duplicateNames = append(duplicateNames, name)
		}
	}
	if len(duplicateNames) > 0 {
		return model.BuildError(
			defaultValidateKey,
			fmt.Errorf("duplicate resource recipe [%s]",
				strings.Join(duplicateNames, ", "),
			),
		)
	}
	return nil
}

func validateAllFrameworks(rcps []*Framework) model.Error {
	nameEncountered := make(map[string]int)
	for _, frameworkRcp := range rcps {
		if err := ValidateFramework(frameworkRcp); err != nil {
			return err
		}
		nameEncountered[frameworkRcp.Name]++
	}
	var duplicateNames []string
	for name, count := range nameEncountered {
		if count > 1 {
			duplicateNames = append(duplicateNames, name)
		}
	}
	if len(duplicateNames) > 0 {
		return model.BuildError(
			defaultValidateKey,
			fmt.Errorf("duplicate framework recipe [%s]",
				strings.Join(duplicateNames, ", "),
			),
		)
	}
	return nil
}

// ValidateResource validates the recipe for a Resource
func ValidateResource(resourceRcp *Resource) model.Error {
	const defaultErrKey = "ValidateResource"
	if err := validator.New().Struct(resourceRcp); err != nil {
		return model.BuildError(defaultErrKey, err)
	}
	return nil
}

// ValidateFramework validates the recipe for a Framework
func ValidateFramework(frameworkRcp *Framework) model.Error {
	const defaultErrKey = "ValidateFramework"
	if err := validator.New().Struct(frameworkRcp); err != nil {
		return model.BuildError(defaultErrKey, err)
	}
	return nil
}
