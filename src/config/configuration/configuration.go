package configuration

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/env"
)

type Command interface {
	WorkingDirectory() *string
	Env() *env.Env

	Validate(c Configuration) error
}

type PrebuiltCommand interface {
	WorkingDirectory() *string
	Env() *env.Env

	Prebuilt() string
	Params() []map[string]string

	Validate(c Configuration) error
}

type ShellCommand struct {
	WorkingDirectory_ *string
	Env_              *env.Env

	Command string
}

func (s ShellCommand) WorkingDirectory() *string {
	return s.WorkingDirectory_
}

func (s ShellCommand) Env() *env.Env {
	return s.Env_
}

func (s ShellCommand) Validate(c Configuration) error {
	return nil // Shell commands are always valid
}

type OperationCommand struct {
	WorkingDirectory_ *string
	Env_              *env.Env

	Operation string
}

func (o OperationCommand) WorkingDirectory() *string {
	return o.WorkingDirectory_
}

func (o OperationCommand) Env() *env.Env {
	return o.Env_
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
