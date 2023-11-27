package config

type Env map[string]string

type YAMLConfig struct {
	Registry *string `yaml:"registry,omitempty" json:"registry,omitempty"`
	Agent    *string `yaml:"agent,omitempty" json:"agent,omitempty"`
}

type YAMLTests map[string]YAMLTest
type YAMLTest struct {
	Directory          *string   `yaml:"directory,omitempty" json:"directory,omitempty"`
	ExcludeDirectories *[]string `yaml:"exclude_directories,omitempty" json:"exclude_directories,omitempty"`
	Description        *string   `yaml:"description,omitempty" json:"description,omitempty"`

	Language  *string `yaml:"language,omitempty" json:"language,omitempty"`
	Framework *string `yaml:"framework,omitempty" json:"framework,omitempty"`

	Run *string `yaml:"run,omitempty" json:"run,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-"`
}

type YAMLImages map[string]YAMLImage
type YAMLImage struct {
	Image *string `yaml:"image,omitempty" json:"image,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-"`
}

type YAMLBuilds map[string]YAMLBuild
type YAMLBuild struct {
	Image     *string `yaml:"image,omitempty" json:"image,omitempty"`
	Directory *string `yaml:"directory,omitempty" json:"directory,omitempty"`
	Build     *string `yaml:"build,omitempty" json:"build,omitempty"`
	Output    *string `yaml:"output,omitempty" json:"output,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-"`
}

type YAMLWorkflowImages map[string]YAMLWorkflowTests
type YAMLWorkflowTests []YAMLWorkflowTest
type YAMLWorkflowTest string
type YAMLWorkflowConditions []string // TODO: More complex conditions
type YAMLWorkflows map[string]YAMLWorkflow
type YAMLWorkflow struct {
	Tests YAMLWorkflowImages `yaml:"tests" json:"tests"`

	Conditions YAMLWorkflowConditions `yaml:"conditions,omitempty" json:"conditions,omitempty"`

	Description *string `yaml:"description,omitempty" json:"description,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty"`

	// Computed
	Name string `yaml:"-" json:"-"`
}

type Config struct {
	Config    YAMLConfig    `yaml:"config" json:"config"`
	Tests     YAMLTests     `yaml:"tests" json:"tests"`
	Images    YAMLImages    `yaml:"images" json:"images"`
	Builds    YAMLBuilds    `yaml:"builds" json:"builds"`
	Workflows YAMLWorkflows `yaml:"workflows" json:"workflows"`
}
