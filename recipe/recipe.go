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
	Path           string   `yaml:"path" validate:"required"`
	FrameworkNames []string `yaml:"framework_names" validate:"required,min=1"`
}

// Framework is a recipe on how and where to read the actual Framework data
type Framework struct {
	Name          string          `yaml:"name" validate:"required"`
	Definitions   []*Definition   `yaml:"definitions"`
	Schemas       []*Schema       `yaml:"schemas"`
	Procedures    []*Procedure    `yaml:"procedures"`
	OutputTargets []*OutputTarget `yaml:"output_targets"`
}

// Definition is a recipe on how and where to read the actual Definition data
type Definition struct {
	Name     string    `yaml:"name" validate:"required"`
	Format   string    `yaml:"format" validate:"required,oneof=json yaml"`
	Type     string    `yaml:"type" validate:"required,oneof=dir file"`
	Path     string    `yaml:"path" validate:"required"`
	Function *Function `yaml:"function"`
}

// Function is a recipe on how to construct a Definition
type Function struct {
	Type string `yaml:"type" validate:"required,oneof=dir file"`
	Path string `yaml:"path" validate:"required"`
}

// Schema is a recipe on how and where to read the actual Schema data
type Schema struct {
	Name string `yaml:"name" validate:"required"`
	Type string `yaml:"type" validate:"required,oneof=dir file"`
	Path string `yaml:"path" validate:"required"`
}

// Procedure is a recipe on how and where to read the actual Procedure data
type Procedure struct {
	Name          string `yaml:"name" validate:"required"`
	Type          string `yaml:"type" validate:"required,oneof=dir file"`
	Path          string `yaml:"path" validate:"required"`
	OutputIsError bool   `yaml:"output_is_error"`
}

// OutputTarget defines how an output is created
type OutputTarget struct {
	Name   string `yaml:"name" validate:"required"`
	Format string `yaml:"format" validate:"required,oneof=json yaml"`
	Type   string `yaml:"type" validate:"required,eq=dir"`
	Path   string `yaml:"path"`
}
