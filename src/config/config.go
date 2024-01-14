package config

import "fmt"

type Env map[string]string

type Command interface {
	WorkingDirectory() *string
	Env() *Env

	Validate(c Configuration) error
}

type PrebuiltCommand interface {
	WorkingDirectory() *string
	Env() *Env

	Prebuilt() string
	Params() []map[string]string

	Validate(c Configuration) error
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

func (s ShellCommand) Validate(c Configuration) error {
	return nil // Shell commands are always valid
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

func (o OperationCommand) Validate(c Configuration) error {
	for _, op := range c.OperationSection {
		if op.Name == o.Operation {
			return ValidateCommands(c, op.Commands)
		}
	}
	return fmt.Errorf("operation '%s' not found", o.Operation)
}

type Test struct {
	Name     string
	Commands []Command

	WorkingDirectory *string
	Env              *Env
}

type TestSection []Test

type Operation struct {
	Name     string
	Commands []Command

	WorkingDirectory *string
	Env              *Env
}

type OperationSection []Operation

type Runtime interface {
	Name() string
	Env() *Env

	Validate(c Configuration) error
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

func (r DockerRuntime) Validate(c Configuration) error {
	return nil // Docker runtimes are always valid
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

func (r BareMetalRuntime) Validate(c Configuration) error {
	return nil // Bare metal runtimes are always valid
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
	OperationSection  OperationSection
	RuntimeSection    RuntimeSection
	BuildSection      BuildSection
	DeploymentSection DeploymentSection
	WorkflowSection   WorkflowSection
	ConfigSection     ConfigSection
}
