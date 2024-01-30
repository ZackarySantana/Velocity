package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zackarysantana/velocity/src/config/configuration"
	"github.com/zackarysantana/velocity/src/env"
	"github.com/zackarysantana/velocity/src/prebuilt"
)

func HydrateConfiguration(raw *RawConfiguration) (*configuration.Configuration, error) {
	var errs []error
	var err error

	config := &configuration.Configuration{}

	config.TestSection, err = HydrateTestSection(raw.TestSection)
	if err != nil {
		errs = append(errs, err)
	}

	config.OperationSection, err = HydrateOperationSection(raw.OperationSection)
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

func HydrateTestSection(raw RawTestSection) (configuration.TestSection, error) {
	var errs []error
	var err error

	testSection := make(configuration.TestSection, len(raw))

	for i, rawTest := range raw {
		testSection[i], err = HydrateTest(rawTest)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return testSection, nil
}

func HydrateTest(raw RawTest) (configuration.Test, error) {
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return configuration.Test{}, err
	}

	test := configuration.Test{
		Name:             raw.Name,
		WorkingDirectory: raw.WorkingDirectory,
		Env:              env,
	}
	test.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return configuration.Test{}, err
	}

	return test, nil
}

func HydrateOperationSection(raw RawOperationSection) (configuration.OperationSection, error) {
	var errs []error
	var err error

	operationSection := make(configuration.OperationSection, len(raw))

	for i, rawOperation := range raw {
		operationSection[i], err = HydrateOperation(rawOperation)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return operationSection, nil
}

func HydrateOperation(raw RawOperation) (configuration.Operation, error) {
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return configuration.Operation{}, err
	}

	operation := configuration.Operation{
		Name:             raw.Name,
		WorkingDirectory: raw.WorkingDirectory,
		Env:              env,
	}
	operation.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return configuration.Operation{}, err
	}

	return operation, nil
}

func HydrateCommands(raw []RawCommand) ([]configuration.Command, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var errs []error
	var err error

	commands := make([]configuration.Command, len(raw))

	for i, rawCommand := range raw {
		commands[i], err = HydrateCommand(rawCommand)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return commands, nil
}

func HydrateCommand(raw RawCommand) (configuration.Command, error) {
	if raw.Prebuilt != nil {
		return HydratePrebuiltCommand(raw)
	} else if raw.Operation != nil {
		return HydrateOperationCommand(raw)
	} else if raw.Command != nil {
		return HydrateShellCommand(raw)
	}

	return nil, fmt.Errorf("invalid command: %v", raw)
}

func HydratePrebuiltCommand(raw RawCommand) (configuration.PrebuiltCommand, error) {
	if raw.Prebuilt == nil {
		return nil, fmt.Errorf("invalid command: %v", raw)
	}
	constructor, err := prebuilt.GetPrebuiltConstructor(*raw.Prebuilt)
	if err != nil {
		return nil, fmt.Errorf("invalid prebuilt command '%s': %w", *raw.Prebuilt, err)
	}
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return nil, err
	}
	params := []map[string]string{}
	if raw.Params != nil {
		params = *raw.Params
	}
	return constructor(
		configuration.CommandInfo{
			WorkingDirectory: raw.WorkingDirectory,
			Env:              env,
		},
		prebuilt.PrebuiltInfo{
			Params: params,
		},
	), err
}

func HydrateOperationCommand(raw RawCommand) (configuration.OperationCommand, error) {
	if raw.Operation == nil {
		return configuration.OperationCommand{}, fmt.Errorf("invalid command: %v", raw)
	}
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return configuration.OperationCommand{}, err
	}
	return configuration.OperationCommand{
		Info: configuration.CommandInfo{
			WorkingDirectory: raw.WorkingDirectory,
			Env:              env,
		},
		Operation: *raw.Operation,
	}, nil
}

func HydrateShellCommand(raw RawCommand) (configuration.ShellCommand, error) {
	if raw.Command == nil {
		return configuration.ShellCommand{}, fmt.Errorf("invalid command: %v", raw)
	}
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return configuration.ShellCommand{}, err
	}
	return configuration.ShellCommand{
		Info: configuration.CommandInfo{
			WorkingDirectory: raw.WorkingDirectory,
			Env:              env,
		},
		Command: *raw.Command,
	}, nil
}

func HydrateRuntimeSection(raw RawRuntimeSection) (configuration.RuntimeSection, error) {
	var errs []error
	var err error

	runtimeSection := make(configuration.RuntimeSection, len(raw))

	for i, rawRuntime := range raw {
		runtimeSection[i], err = HydrateRuntime(rawRuntime)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return runtimeSection, nil
}

func HydrateRuntime(raw RawRuntime) (configuration.Runtime, error) {
	env, err := HydrateEnv(raw.Env)
	if err != nil {
		return nil, err
	}

	if raw.Image != nil {
		return configuration.DockerRuntime{
			Name_: raw.Name,
			Env_:  env,

			Image: *raw.Image,
		}, nil
	} else if raw.Machine != nil {
		return configuration.BareMetalRuntime{
			Name_: raw.Name,
			Env_:  env,

			Machine: raw.Machine,
		}, nil
	}

	return nil, fmt.Errorf("invalid runtime: %v", raw)
}

func HydrateBuildSection(raw RawBuildSection) (configuration.BuildSection, error) {
	var errs []error
	var err error

	buildSection := make(configuration.BuildSection, len(raw))

	for i, rawBuild := range raw {
		buildSection[i], err = HydrateBuild(rawBuild)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return buildSection, nil
}

func HydrateBuild(raw RawBuild) (configuration.Build, error) {
	var err error

	build := configuration.Build{
		Name:         raw.Name,
		BuildRuntime: raw.BuildRuntime,
		Output:       raw.Output,

		OutputRuntime: raw.OutputRuntime,
		OutputCmd:     raw.OutputCmd,
	}

	build.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return configuration.Build{}, err
	}

	return build, nil
}

func HydrateDeploymentSection(raw RawDeploymentSection) (configuration.DeploymentSection, error) {
	var errs []error
	var err error

	deploymentSection := make(configuration.DeploymentSection, len(raw))

	for i, rawDeployment := range raw {
		deploymentSection[i], err = HydrateDeployment(rawDeployment)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return deploymentSection, nil
}

func HydrateDeployment(raw RawDeployment) (configuration.Deployment, error) {
	var err error

	deployment := configuration.Deployment{
		Name:      raw.Name,
		Workflows: raw.Workflows,
	}

	deployment.Commands, err = HydrateCommands(raw.Commands)
	if err != nil {
		return configuration.Deployment{}, err
	}

	return deployment, nil
}

func HydrateWorkflowSection(raw RawWorkflowSection) (configuration.WorkflowSection, error) {
	var errs []error
	var err error

	workflowSection := make(configuration.WorkflowSection, len(raw))

	for i, rawWorkflow := range raw {
		workflowSection[i], err = HydrateWorkflow(rawWorkflow)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return workflowSection, nil
}

func HydrateWorkflow(raw RawWorkflow) (configuration.Workflow, error) {
	var err error

	workflow := configuration.Workflow{
		Name: raw.Name,
	}

	workflow.Groups, err = HydrateWorkflowGroups(raw.Groups)
	if err != nil {
		return configuration.Workflow{}, err
	}

	return workflow, nil
}

func HydrateWorkflowGroups(raw []RawWorkflowGroup) ([]configuration.WorkflowGroup, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var errs []error
	var err error

	workflowGroups := make([]configuration.WorkflowGroup, len(raw))

	for i, rawWorkflowGroup := range raw {
		workflowGroups[i], err = HydrateWorkflowGroup(rawWorkflowGroup)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return workflowGroups, nil
}

func HydrateWorkflowGroup(raw RawWorkflowGroup) (configuration.WorkflowGroup, error) {
	return configuration.WorkflowGroup(raw), nil
}

func HydrateConfigSection(raw RawConfigSection) (configuration.ConfigSection, error) {
	return configuration.ConfigSection(raw), nil
}

func HydrateEnv(raw *RawEnv) (*env.Env, error) {
	if raw == nil {
		return nil, nil
	}
	var errs []error
	envs := make(env.Env, len(*raw))

	for _, rawEnv := range *raw {
		name, value, err := HydrateEnvLine(rawEnv)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		envs[name] = value
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
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
