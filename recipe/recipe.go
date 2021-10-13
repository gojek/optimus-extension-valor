package recipe

// Recipe is the main structure for storing recipe execution flow
type Recipe struct {
	Frameworks []*Framework `yaml:"frameworks" validate:"required"`
	Resources  []*Resource  `yaml:"resources"`
}

// Framework is a recipe on how how and where to read the actual Framework data
type Framework struct {
	AllowError    bool        `yaml:"allow_error"`
	Name          string      `yaml:"name" validate:"required"`
	Schemas       []*Metadata `yaml:"schemas"`
	Definitions   []*Metadata `yaml:"definitions"`
	Procedures    []*Metadata `yaml:"procedures"`
	OutputTargets []*Metadata `yaml:"output_targets"`
}

// Metadata holds information to where and how a data is stored
type Metadata struct {
	Name          string `yaml:"name" validate:"required"`
	Format        string `yaml:"format" validate:"required"`
	Type          string `yaml:"type" validate:"required"`
	Path          string `yaml:"path" validate:"required"`
	OutputIsError bool   `yaml:"output_is_error"`
}

// Resource is a recipe on how how and where to read the actual Resource data
type Resource struct {
	Name           string   `yaml:"name" validate:"required"`
	Format         string   `yaml:"format" validate:"required"`
	Type           string   `yaml:"type" validate:"required"`
	Path           string   `yaml:"path" validate:"required"`
	FrameworkNames []string `yaml:"framework_names" validate:"required"`
}
