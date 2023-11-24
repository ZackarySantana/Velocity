package config

type Env map[string]string

type YAMLConfig struct {
	Registry *string `yaml:"registry,omitempty" json:"registry,omitempty"`
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
}

type YAMLImages map[string]YAMLImage
type YAMLImage struct {
	Image *string `yaml:"image,omitempty" json:"image,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
}

type YAMLWorkflowTests map[string]YAMLWorkflowTest
type YAMLWorkflowTest []string
type YAMLWorkflowConditions []string // TODO: More complex conditions
type YAMLWorkflows map[string]YAMLWorkflow
type YAMLWorkflow struct {
	Tests YAMLWorkflowTests `yaml:"tests" json:"tests"`

	Conditions YAMLWorkflowConditions `yaml:"conditions,omitempty" json:"conditions,omitempty"`

	Description *string `yaml:"description,omitempty" json:"description,omitempty"`

	Env *map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
}

type Config struct {
	Config    YAMLConfig    `yaml:"config" json:"config"`
	Tests     YAMLTests     `yaml:"tests" json:"tests"`
	Images    YAMLImages    `yaml:"images" json:"images"`
	Workflows YAMLWorkflows `yaml:"workflows" json:"workflows"`
}
