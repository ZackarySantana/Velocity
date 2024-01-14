package config

import (
	"gopkg.in/yaml.v3"
)

type RawEnv []string

type RawCommand struct {
	WorkingDirectory *string
	Env              *RawEnv

	// Prebuilt
	Prebuilt *string
	Params   *[]map[string]string

	// Command
	Command *string

	// Operations
	Operation *string
}

type RawTest struct {
	Name     string
	Commands []RawCommand

	WorkingDirectory *string
	Env              *RawEnv
}

type RawTestSection []RawTest

type RawRuntime struct {
	Name string

	Env *RawEnv

	// Docker runtime
	Image *string

	// Bare metal
	Machine *string
}

type RawRuntimeSection []RawRuntime

type RawBuild struct {
	Name         string
	BuildRuntime string
	Output       string
	Commands     []RawCommand

	OutputRuntime *string
	OutputCmd     *string
}

type RawBuildSection []RawBuild

type RawDeployment struct {
	Name     string
	Commands []RawCommand

	Workflows []string
}

type RawDeploymentSection []RawDeployment

type RawWorkflowGroup struct {
	Name string

	Runtimes []string

	Tests []string
}

type RawWorkflow struct {
	Name string

	Groups []RawWorkflowGroup
}

type RawWorkflowSection []RawWorkflow

type RawConfigSection struct {
	Project string

	Server *string
	UI     *string
}

type RawConfiguration struct {
	TestSection       RawTestSection
	RuntimeSection    RawRuntimeSection
	BuildSection      RawBuildSection
	DeploymentSection RawDeploymentSection
	WorkflowSection   RawWorkflowSection
	ConfigSection     RawConfigSection
}

func Parse(data []byte) (*RawConfiguration, error) {
	var config RawConfiguration
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
