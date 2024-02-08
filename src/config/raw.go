package config

import (
	"gopkg.in/yaml.v3"
)

type RawEnv []string

type RawCommand struct {
	WorkingDirectory *string `yaml:"working_dir" json:"working_dir" bson:"working_dir"`
	Env              *RawEnv `yaml:"env" json:"env" bson:"env"`

	// Prebuilt
	Prebuilt *string              `yaml:"prebuilt" json:"prebuilt" bson:"prebuilt"`
	Params   *[]map[string]string `yaml:"params" json:"params" bson:"params"`

	// Command
	Command *string `yaml:"command" json:"command" bson:"command"`

	// Operations
	Operation *string `yaml:"operation" json:"operation" bson:"operation"`
}

type RawTest struct {
	Name     string       `yaml:"name" json:"name" bson:"name"`
	Commands []RawCommand `yaml:"commands" json:"commands" bson:"commands"`

	WorkingDirectory *string `yaml:"working_dir" json:"working_dir" bson:"working_dir"`
	Env              *RawEnv `yaml:"env" json:"env" bson:"env"`
}

type RawTestSection []RawTest

type RawOperation struct {
	Name     string       `yaml:"name" json:"name" bson:"name"`
	Commands []RawCommand `yaml:"commands" json:"commands" bson:"commands"`

	WorkingDirectory *string `yaml:"working_dir" json:"working_dir" bson:"working_dir"`
	Env              *RawEnv `yaml:"env" json:"env" bson:"env"`
}

type RawOperationSection []RawOperation

type RawRuntime struct {
	Name string `yaml:"name" json:"name" bson:"name"`

	Env *RawEnv `yaml:"env" json:"env" bson:"env"`

	// Docker runtime
	Image *string `yaml:"image" json:"image" bson:"image"`

	// Bare metal
	Machine *string `yaml:"machine" json:"machine" bson:"machine"`
}

type RawRuntimeSection []RawRuntime

type RawBuild struct {
	Name         string       `yaml:"name" json:"name" bson:"name"`
	BuildRuntime string       `yaml:"build_runtime" json:"build_runtime" bson:"build_runtime"`
	Output       string       `yaml:"output" json:"output" bson:"output"`
	Commands     []RawCommand `yaml:"commands" json:"commands" bson:"commands"`

	OutputRuntime *string `yaml:"output_runtime" json:"output_runtime" bson:"output_runtime"`
	OutputCmd     *string `yaml:"output_cmd" json:"output_cmd" bson:"output_cmd"`
}

type RawBuildSection []RawBuild

type RawDeployment struct {
	Name     string       `yaml:"name" json:"name" bson:"name"`
	Commands []RawCommand `yaml:"commands" json:"commands" bson:"commands"`

	Workflows []string `yaml:"workflows" json:"workflows" bson:"workflows"`
}

type RawDeploymentSection []RawDeployment

type RawWorkflowGroup struct {
	Name string `yaml:"name" json:"name" bson:"name"`

	Runtimes []string `yaml:"runtimes" json:"runtimes" bson:"runtimes"`

	Tests []string `yaml:"tests" json:"tests" bson:"tests"`
}

type RawWorkflow struct {
	Name string `yaml:"name" json:"name" bson:"name"`

	Groups []RawWorkflowGroup `yaml:"groups" json:"groups" bson:"groups"`
}

type RawWorkflowSection []RawWorkflow

type RawConfigSection struct {
	Project string `yaml:"project" json:"project" bson:"project"`

	Server *string `yaml:"server" json:"server" bson:"server"`
	UI     *string `yaml:"ui" json:"ui" bson:"ui"`
}

type RawConfiguration struct {
	TestSection       RawTestSection       `yaml:"tests" json:"tests" bson:"tests"`
	OperationSection  RawOperationSection  `yaml:"operations" json:"operations" bson:"operations"`
	RuntimeSection    RawRuntimeSection    `yaml:"runtimes" json:"runtimes" bson:"runtimes"`
	BuildSection      RawBuildSection      `yaml:"builds" json:"builds" bson:"builds"`
	DeploymentSection RawDeploymentSection `yaml:"deployments" json:"deployments" bson:"deployments"`
	WorkflowSection   RawWorkflowSection   `yaml:"workflows" json:"workflows" bson:"workflows"`
	ConfigSection     RawConfigSection     `yaml:"config" json:"config" bson:"config"`
}

func Parse(data []byte) (*RawConfiguration, error) {
	var config RawConfiguration
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
