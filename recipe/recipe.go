package recipe

// Recipe is the main structure for storing recipe execution flow
type Recipe struct {
	Resources  []*Resource  `yaml:"resources" validate:"required"`
	Frameworks []*Framework `yaml:"frameworks" validate:"required"`
}

// Resource is a recipe on how and where to read the actual Resource data
type Resource struct {
	Name           string   `yaml:"name" validate:"required"`
	Format         string   `yaml:"format" validate:"required,oneof=json yaml"`
	Type           string   `yaml:"type" validate:"required,oneof=dir file"`
	RegexPattern   string   `yaml:"regex_pattern"`
	Path           string   `yaml:"path" validate:"required"`
	BatchSize      int      `yaml:"batch_size"`
	FrameworkNames []string `yaml:"framework_names" validate:"required,min=1"`
}

// Framework is a recipe on how and where to read the actual Framework data
type Framework struct {
	Name        string        `yaml:"name" validate:"required"`
	Schemas     []*Schema     `yaml:"schemas"`
	Definitions []*Definition `yaml:"definitions"`
	Procedures  []*Procedure  `yaml:"procedures"`
}

// Definition is a recipe on how and where to read the actual Definition data
type Definition struct {
	Name         string    `yaml:"name" validate:"required"`
	Format       string    `yaml:"format" validate:"required,oneof=json yaml"`
	Type         string    `yaml:"type" validate:"required,oneof=dir file"`
	Path         string    `yaml:"path" validate:"required"`
	RegexPattern string    `yaml:"regex_pattern"`
	Function     *Function `yaml:"function"`
}

// Function is a recipe on how to construct a Definition
type Function struct {
	Type string `yaml:"type" validate:"required,oneof=dir file"`
	Path string `yaml:"path" validate:"required"`
}

// Schema is a recipe on how and where to read the actual Schema data
type Schema struct {
	Name   string  `yaml:"name" validate:"required"`
	Type   string  `yaml:"type" validate:"required,oneof=dir file"`
	Path   string  `yaml:"path" validate:"required"`
	Output *Output `yaml:"output"`
}

// Procedure is a recipe on how and where to read the actual Procedure data
type Procedure struct {
	Name   string  `yaml:"name" validate:"required"`
	Type   string  `yaml:"type" validate:"required,oneof=dir file"`
	Path   string  `yaml:"path" validate:"required"`
	Output *Output `yaml:"output"`
}

// Output defines how the last procedure output is written
type Output struct {
	TreatAs string    `yaml:"treat_as" validate:"required,oneof=info warning error success"`
	Targets []*Target `yaml:"targets" validate:"required,min=1"`
}

// Target defines how an output is written to the targetted stream
type Target struct {
	Name   string `yaml:"name" validate:"required"`
	Format string `yaml:"format" validate:"required,oneof=json yaml"`
	Type   string `yaml:"type" validate:"required,eq=dir"`
	Path   string `yaml:"path"`
}
