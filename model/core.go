package model

const (
	// SkipEmptyValue is an empty value that is not considered as value
	SkipEmptyValue = ""
	// SkipNullValue is a null value that is not considered as value
	SkipNullValue = "null\n"
)

// IsSkipResult maps skip values
var IsSkipResult = map[string]bool{
	SkipEmptyValue: true,
	SkipNullValue:  true,
}

// Data contains data information
type Data struct {
	Type    string
	Path    string
	Content []byte
}

// Resource contains resource data to be processed
type Resource struct {
	Name           string
	FrameworkNames []string
	ListOfData     []*Data
}

// Framework contains information on how to process a Resource
type Framework struct {
	Name          string
	Definitions   []*Definition
	Schemas       []*Schema
	Procedures    []*Procedure
	OutputTargets []*OutputTarget
}

// Definition is the definition that could be used in a Procedure
type Definition struct {
	Name         string
	ListOfData   []*Data
	FunctionData *Data
}

// Schema contains information on Schema information defined by the user
type Schema struct {
	Name string
	Data *Data
}

// Procedure contains information on Procedure information defined by the user
type Procedure struct {
	Name          string
	OutputIsError bool
	Data          *Data
}

// OutputTarget describes the target of output for each Framework
type OutputTarget struct {
	Name   string
	Format string
	Type   string
	Path   string
}
