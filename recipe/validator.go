package recipe

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate validates the recipe
func Validate(rcp *Recipe) error {
	if err := validator.New().Struct(rcp); err != nil {
		return err
	}
	if err := validateAllResources(rcp.Resources); err != nil {
		return err
	}
	return validateAllFrameworks(rcp.Frameworks)
}

func validateAllResources(rcps []*Resource) error {
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
		return fmt.Errorf("duplicate resource recipe [%s]",
			strings.Join(duplicateNames, ", "))
	}
	return nil
}

func validateAllFrameworks(rcps []*Framework) error {
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
		return fmt.Errorf("duplicate framework recipe [%s]",
			strings.Join(duplicateNames, ", "))
	}
	return nil
}

// ValidateResource validates the recipe for a Resource
func ValidateResource(resourceRcp *Resource) error {
	if err := validator.New().Struct(resourceRcp); err != nil {
		return err
	}
	return nil
}

// ValidateFramework validates the recipe for a Framework
func ValidateFramework(frameworkRcp *Framework) error {
	if err := validator.New().Struct(frameworkRcp); err != nil {
		return err
	}
	return nil
}
