package config

type Env map[string]string

type YAMLConfig struct {
	Registry *string `yaml:"registry,omitempty" json:"registry,omitempty" bson:"registry,omitempty"`
	Server   *string `yaml:"server,omitempty" json:"server,omitempty" bson:"server,omitempty"`
}

type YAMLTests map[string]YAMLTest
type YAMLTest struct {
	Directory          *string   `yaml:"directory,omitempty" json:"directory,omitempty" bson:"directory,omitempty"`
	ExcludeDirectories *[]string `yaml:"exclude_directories,omitempty" json:"exclude_directories,omitempty" bson:"exclude_directories,omitempty"`
	Description        *string   `yaml:"description,omitempty" json:"description,omitempty" bson:"description,omitempty"`

	Language  *string `yaml:"language,omitempty" json:"language,omitempty" bson:"language,omitempty"`
	Framework *string `yaml:"framework,omitempty" json:"framework,omitempty" bson:"framework,omitempty"`

	Run *string `yaml:"run,omitempty" json:"run,omitempty" bson:"run,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty" bson:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-" bson:"-"`
}

type YAMLImages map[string]YAMLImage
type YAMLImage struct {
	Image *string `yaml:"image,omitempty" json:"image,omitempty" bson:"image,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty" bson:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-" bson:"-"`
}

type YAMLBuilds map[string]YAMLBuild
type YAMLBuild struct {
	Image     *string `yaml:"image,omitempty" json:"image,omitempty" bson:"image,omitempty"`
	Directory *string `yaml:"directory,omitempty" json:"directory,omitempty" bson:"directory,omitempty"`
	Build     *string `yaml:"build,omitempty" json:"build,omitempty" bson:"build,omitempty"`
	Output    *string `yaml:"output,omitempty" json:"output,omitempty" bson:"output,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty" bson:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-" bson:"-"`
}

type YAMLWorkflowImages map[string]YAMLWorkflowTests
type YAMLWorkflowTests []YAMLWorkflowTest
type YAMLWorkflowTest string
type YAMLWorkflowConditions []string // TODO: More complex conditions
type YAMLWorkflows map[string]YAMLWorkflow
type YAMLWorkflow struct {
	Tests YAMLWorkflowImages `yaml:"tests" json:"tests" bson:"tests"`

	Conditions YAMLWorkflowConditions `yaml:"conditions,omitempty" json:"conditions,omitempty" bson:"conditions,omitempty"`

	Description *string `yaml:"description,omitempty" json:"description,omitempty" bson:"description,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty" bson:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-" bson:"-"`
}

type Config struct {
	Config    YAMLConfig    `yaml:"config" json:"config" bson:"config"`
	Tests     YAMLTests     `yaml:"tests" json:"tests" bson:"tests"`
	Images    YAMLImages    `yaml:"images" json:"images" bson:"images"`
	Builds    YAMLBuilds    `yaml:"builds" json:"builds" bson:"builds"`
	Workflows YAMLWorkflows `yaml:"workflows" json:"workflows" bson:"workflows"`

	// Computed
	Path string `yaml:"-" json:"-" bson:"-"`
}
