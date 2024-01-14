package config

import (
	"gopkg.in/yaml.v3"
)

type RawEnv []string

type RawCommand struct {
	WorkingDirectory *string `yaml:"working_dir"`
	Env              *RawEnv `yaml:"env"`

	// Prebuilt
	Prebuilt *string              `yaml:"prebuilt"`
	Params   *[]map[string]string `yaml:"params"`

	// Command
	Command *string `yaml:"command"`

	// Operations
	Operation *string `yaml:"operation"`
}

type RawTest struct {
	Name     string       `yaml:"name"`
	Commands []RawCommand `yaml:"commands"`

	WorkingDirectory *string `yaml:"working_dir"`
	Env              *RawEnv `yaml:"env"`
}

type RawTestSection []RawTest

type RawRuntime struct {
	Name string `yaml:"name"`

	Env *RawEnv `yaml:"env"`

	// Docker runtime
	Image *string `yaml:"image"`

	// Bare metal
	Machine *string `yaml:"machine"`
}

type RawRuntimeSection []RawRuntime

type RawBuild struct {
	Name         string       `yaml:"name"`
	BuildRuntime string       `yaml:"build_runtime"`
	Output       string       `yaml:"output"`
	Commands     []RawCommand `yaml:"commands"`

	OutputRuntime *string `yaml:"output_runtime"`
	OutputCmd     *string `yaml:"output_cmd"`
}

type RawBuildSection []RawBuild

type RawDeployment struct {
	Name     string       `yaml:"name"`
	Commands []RawCommand `yaml:"commands"`

	Workflows []string `yaml:"workflows"`
}

type RawDeploymentSection []RawDeployment

type RawWorkflowGroup struct {
	Name string `yaml:"name"`

	Runtimes []string `yaml:"runtimes"`

	Tests []string `yaml:"tests"`
}

type RawWorkflow struct {
	Name string `yaml:"name"`

	Groups []RawWorkflowGroup `yaml:"groups"`
}

type RawWorkflowSection []RawWorkflow

type RawConfigSection struct {
	Project string `yaml:"project"`

	Server *string `yaml:"server"`
	UI     *string `yaml:"ui"`
}

type RawConfiguration struct {
	TestSection       RawTestSection       `yaml:"tests"`
	RuntimeSection    RawRuntimeSection    `yaml:"runtimes"`
	BuildSection      RawBuildSection      `yaml:"builds"`
	DeploymentSection RawDeploymentSection `yaml:"deployments"`
	WorkflowSection   RawWorkflowSection   `yaml:"workflows"`
	ConfigSection     RawConfigSection     `yaml:"config"`
}

func Parse(data []byte) (*RawConfiguration, error) {
	var config RawConfiguration
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
