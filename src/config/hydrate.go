package config

import (
	"errors"
	"fmt"
	"strings"
)

type Env map[string]string

type Command interface {
	WorkingDirectory() *string
	Env() *Env
}

type PrebuiltCommand interface {
	WorkingDirectory() *string
	Env() *Env

	Prebuilt() string
	Params() []map[string]string
}

type ShellCommand struct {
	WorkingDirectory_ *string
	Env_              *Env

	Command string
}

func (s ShellCommand) WorkingDirectory() *string {
	return s.WorkingDirectory_
}

func (s ShellCommand) Env() *Env {
	return s.Env_
}

type OperationCommand struct {
	WorkingDirectory_ *string
	Env_              *Env

	Operation string
}

func (o OperationCommand) WorkingDirectory() *string {
	return o.WorkingDirectory_
}

func (o OperationCommand) Env() *Env {
	return o.Env_
}

type Test struct {
	Name     string
	Commands []Command

	WorkingDirectory *string
	Env              *Env
}

type TestSection []Test

type Runtime interface {
	Name() string
	Env() *Env
}

type DockerRuntime struct {
	Name_ string
	Env_  *Env

	Image string
}

func (r DockerRuntime) Name() string {
	return r.Name_
}

func (r DockerRuntime) Env() *Env {
	return r.Env_
}

type BareMetalRuntime struct {
	Name_ string
	Env_  *Env

	Machine *string
}

func (r BareMetalRuntime) Name() string {
	return r.Name_
}

func (r BareMetalRuntime) Env() *Env {
	return r.Env_
}

type RuntimeSection []Runtime

type Build struct {
	Name         string
	BuildRuntime string
	Output       string
	Commands     []Command

	OutputRuntime *string
	OutputCmd     *string
}

type BuildSection []Build

type Deployment struct {
	Name     string
	Commands []Command

	Workflows []string
}

type DeploymentSection []Deployment

type WorkflowGroup struct {
	Name string

	Runtimes []string

	Tests []string
}

type Workflow struct {
	Name string

	Groups []WorkflowGroup
}

type WorkflowSection []Workflow

type ConfigSection struct {
	Project string

	Server *string
	UI     *string
}

type Configuration struct {
	TestSection       TestSection
	RuntimeSection    RuntimeSection
	BuildSection      BuildSection
	DeploymentSection DeploymentSection
	WorkflowSection   WorkflowSection
	ConfigSection     ConfigSection
}

func HydrateConfiguration(raw *RawConfiguration) (*Configuration, error) {
	var errs []error
	var err error

	config := &Configuration{}

	config.TestSection, err = HydrateTestSection(raw.TestSection)
	if err != nil {
		errs = append(errs, err)
	}

	config.RuntimeSection, err = HydrateRuntimeSection(raw.RuntimeSection)
	if err != nil {
		errs = append(errs, err)
	}

	config.BuildSection, err = HydrateBuildSection(raw.BuildSection)
	if err != nil {
		errs = append(errs, err)
	}

	config.DeploymentSection, err = HydrateDeploymentSection(raw.DeploymentSection)
	if err != nil {
		errs = append(errs, err)
	}

	config.WorkflowSection, err = HydrateWorkflowSection(raw.WorkflowSection)
	if err != nil {
		errs = append(errs, err)
	}

	config.ConfigSection, err = HydrateConfigSection(raw.ConfigSection)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return config, nil
}

func HydrateTestSection(raw RawTestSection) (TestSection, error) {
	var errs []error
	var err error

	testSection := make(TestSection, len(raw))

	for i, rawTest := range raw {
		testSection[i], err = HydrateTest(rawTest)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return testSection, nil
}

func HydrateTest(raw RawTest) (Test, error) {
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return Test{}, err
	}

	test := Test{
		Name:             raw.Name,
		WorkingDirectory: raw.WorkingDirectory,
		Env:              env,
	}
	test.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return Test{}, err
	}

	return test, nil
}

func HydrateCommands(raw []RawCommand) ([]Command, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var errs []error
	var err error

	commands := make([]Command, len(raw))

	for i, rawCommand := range raw {
		commands[i], err = HydrateCommand(rawCommand)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return commands, nil
}

func HydrateCommand(raw RawCommand) (Command, error) {
	if raw.Prebuilt != nil {
		return HydratePrebuiltCommand(raw)
	} else if raw.Operation != nil {
		return HydrateOperationCommand(raw)
	} else if raw.Command != nil {
		return HydrateShellCommand(raw)
	}

	return nil, fmt.Errorf("invalid command: %v", raw)
}

func HydratePrebuiltCommand(raw RawCommand) (PrebuiltCommand, error) {
	return nil, nil // TODO: implement
}

func HydrateOperationCommand(raw RawCommand) (OperationCommand, error) {
	if raw.Operation == nil {
		return OperationCommand{}, fmt.Errorf("invalid command: %v", raw)
	}
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return OperationCommand{}, err
	}
	return OperationCommand{
		WorkingDirectory_: raw.WorkingDirectory,
		Env_:              env,
		Operation:         *raw.Operation,
	}, nil
}

func HydrateShellCommand(raw RawCommand) (ShellCommand, error) {
	if raw.Command == nil {
		return ShellCommand{}, fmt.Errorf("invalid command: %v", raw)
	}
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return ShellCommand{}, err
	}
	return ShellCommand{
		WorkingDirectory_: raw.WorkingDirectory,
		Env_:              env,
		Command:           *raw.Command,
	}, nil
}

func HydrateRuntimeSection(raw RawRuntimeSection) (RuntimeSection, error) {
	var errs []error
	var err error

	runtimeSection := make(RuntimeSection, len(raw))

	for i, rawRuntime := range raw {
		runtimeSection[i], err = HydrateRuntime(rawRuntime)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return runtimeSection, nil
}

func HydrateRuntime(raw RawRuntime) (Runtime, error) {
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return nil, err
	}

	if raw.Image != nil {
		return DockerRuntime{
			Name_: raw.Name,
			Env_:  env,

			Image: *raw.Image,
		}, nil
	} else if raw.Machine != nil {
		return BareMetalRuntime{
			Name_: raw.Name,
			Env_:  env,

			Machine: raw.Machine,
		}, nil
	}

	return nil, fmt.Errorf("invalid runtime: %v", raw)
}

func HydrateBuildSection(raw RawBuildSection) (BuildSection, error) {
	var errs []error
	var err error

	buildSection := make(BuildSection, len(raw))

	for i, rawBuild := range raw {
		buildSection[i], err = HydrateBuild(rawBuild)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return buildSection, nil
}

func HydrateBuild(raw RawBuild) (Build, error) {
	var err error

	build := Build{
		Name:         raw.Name,
		BuildRuntime: raw.BuildRuntime,
		Output:       raw.Output,

		OutputRuntime: raw.OutputRuntime,
		OutputCmd:     raw.OutputCmd,
	}

	build.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return Build{}, err
	}

	return build, nil
}

func HydrateDeploymentSection(raw RawDeploymentSection) (DeploymentSection, error) {
	var errs []error
	var err error

	deploymentSection := make(DeploymentSection, len(raw))

	for i, rawDeployment := range raw {
		deploymentSection[i], err = HydrateDeployment(rawDeployment)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return deploymentSection, nil
}

func HydrateDeployment(raw RawDeployment) (Deployment, error) {
	var err error

	deployment := Deployment{
		Name:      raw.Name,
		Workflows: raw.Workflows,
	}

	deployment.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return Deployment{}, err
	}

	return deployment, nil
}

func HydrateWorkflowSection(raw RawWorkflowSection) (WorkflowSection, error) {
	var errs []error
	var err error

	workflowSection := make(WorkflowSection, len(raw))

	for i, rawWorkflow := range raw {
		workflowSection[i], err = HydrateWorkflow(rawWorkflow)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return workflowSection, nil
}

func HydrateWorkflow(raw RawWorkflow) (Workflow, error) {
	var err error

	workflow := Workflow{
		Name: raw.Name,
	}

	workflow.Groups, err = HydrateWorkflowGroups(raw.Groups)
	if err != nil {
		return Workflow{}, err
	}

	return workflow, nil
}

func HydrateWorkflowGroups(raw []RawWorkflowGroup) ([]WorkflowGroup, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var errs []error
	var err error

	workflowGroups := make([]WorkflowGroup, len(raw))

	for i, rawWorkflowGroup := range raw {
		workflowGroups[i], err = HydrateWorkflowGroup(rawWorkflowGroup)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return workflowGroups, nil
}

func HydrateWorkflowGroup(raw RawWorkflowGroup) (WorkflowGroup, error) {
	return WorkflowGroup(raw), nil
}

func HydrateConfigSection(raw RawConfigSection) (ConfigSection, error) {
	return ConfigSection(raw), nil
}

func HydrateEnv(raw *RawEnv) (*Env, error) {
	if raw == nil {
		return nil, nil
	}
	var errs []error
	envs := make(Env, len(*raw))

	for _, rawEnv := range *raw {
		name, value, err := HydrateEnvLine(rawEnv)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		envs[name] = value
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &envs, nil
}

func HydrateEnvLine(raw string) (string, string, error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid env line: %s", raw)
	}

	envName := parts[0]
	envValue := parts[1]
	if len(envName) == 0 {
		return "", "", fmt.Errorf("invalid env line: %s", raw)
	}
	if len(envValue) == 0 {
		return envName, "", nil
	}

	if envValue[0] == '"' {
		if envValue[len(envValue)-1] != '"' {
			return "", "", fmt.Errorf("invalid env line: %s", raw)
		}
	}

	if envValue[0] == '\'' {
		if envValue[len(envValue)-1] != '\'' {
			return "", "", fmt.Errorf("invalid env line: %s", raw)
		}
	}

	if (envValue[0] == '"' && envValue[len(envValue)-1] == '"') ||
		(envValue[0] == '\'' && envValue[len(envValue)-1] == '\'') {
		envValue = envValue[1 : len(envValue)-1]
	}

	return envName, envValue, nil
}
