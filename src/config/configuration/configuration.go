package configuration

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/env"
)

type CommandInfo struct {
	WorkingDirectory *string
	Env              *env.Env

	Params []map[string]string
}

type Command interface {
	GetInfo() CommandInfo
	Validate(c Configuration) error
}

type PrebuiltCommand interface {
	GetInfo() CommandInfo
	Validate(c Configuration) error

	Prebuilt() string
}

type ShellCommand struct {
	Info CommandInfo

	Command string
}

func (s ShellCommand) GetInfo() CommandInfo {
	return s.Info
}

func (s ShellCommand) Validate(c Configuration) error {
	return nil // Shell commands are always valid
}

type OperationCommand struct {
	Info CommandInfo

	Operation string
}

func (o OperationCommand) GetInfo() CommandInfo {
	return o.Info
}

func (o OperationCommand) Validate(c Configuration) error {
	for _, op := range c.OperationSection {
		if op.Name == o.Operation {
			return nil
		}
	}
	return fmt.Errorf("operation '%s' not found", o.Operation)
}

type Test struct {
	Name     string
	Commands []Command

	WorkingDirectory *string
	Env              *env.Env
}

type TestSection []Test

type Operation struct {
	Name     string
	Commands []Command

	WorkingDirectory *string
	Env              *env.Env
}

type OperationSection []Operation

type Runtime interface {
	Name() string
	Env() *env.Env

	Validate(c Configuration) error
}

type DockerRuntime struct {
	Name_ string
	Env_  *env.Env

	Image string
}

func (r DockerRuntime) Name() string {
	return r.Name_
}

func (r DockerRuntime) Env() *env.Env {
	return r.Env_
}

func (r DockerRuntime) Validate(c Configuration) error {
	return nil // Docker runtimes are always valid
}

type BareMetalRuntime struct {
	Name_ string
	Env_  *env.Env

	Machine *string
}

func (r BareMetalRuntime) Name() string {
	return r.Name_
}

func (r BareMetalRuntime) Env() *env.Env {
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
	Tests    []string
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
