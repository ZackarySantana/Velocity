package config

import (
	"errors"
	"fmt"

	"github.com/zackarysantana/velocity/internal/utils/errors2"
	"github.com/zackarysantana/velocity/internal/utils/slices"
)

func ValidateConfiguration(c Configuration) error {
	var errs []error

	errs = append(errs, ValidateTestSectionPartial(c))
	errs = append(errs, ValidateOperationSectionPartial(c))
	errs = append(errs, ValidateRuntimeSectionPartial(c))
	errs = append(errs, ValidateBuildSectionPartial(c))
	errs = append(errs, ValidateDeploymentSectionPartial(c))
	errs = append(errs, ValidateWorkflowSectionPartial(c))
	errs = append(errs, ValidateConfigSectionPartial(c))

	return errors2.JoinWithHead(errors.New("validating configuration"), errs...)
}

func ValidateTestSectionPartial(c Configuration) error {
	var errs []error

	for _, t := range c.TestSection {
		errs = append(errs, ValidateTestPartial(c, t))
	}

	return errors2.JoinWithHead(errors.New("test section"), errs...)
}

func ValidateTestPartial(c Configuration, t Test) error {
	var errs []error

	for _, cmd := range t.Commands {
		errs = append(errs, ValidateCommandPartial(c, cmd))
	}

	if err := errors2.JoinWithHead(fmt.Errorf("validating test '%s'", t.Name), errs...); err != nil {
		return err
	}

	return nil
}

func ValidateCommandsPartial(c Configuration, cmds []Command) error {
	var errs []error
	for _, cmd := range cmds {
		if err := ValidateCommandPartial(c, cmd); err != nil {
			errs = append(errs, err)
		}
	}
	return errors2.JoinWithHead(errors.New("commands have errors"), errs...)
}

func ValidateCommandPartial(c Configuration, cmd Command) error {
	if cmd == nil {
		return fmt.Errorf("command is nil")
	}
	return cmd.Validate(c)
}

func ValidateOperationSectionPartial(c Configuration) error {
	var errs []error

	for _, o := range c.OperationSection {
		errs = append(errs, ValidateOperationPartial(c, o))
	}

	return errors2.JoinWithHead(errors.New("operation section"), errs...)
}

func ValidateOperationPartial(c Configuration, o Operation) error {
	var errs []error

	for _, cmd := range o.Commands {
		errs = append(errs, ValidateCommandPartial(c, cmd))
	}

	if err := errors2.JoinWithHead(fmt.Errorf("validating operation '%s'", o.Name), errs...); err != nil {
		return err
	}

	return nil
}

func ValidateRuntimeSectionPartial(c Configuration) error {
	var errs []error

	for _, r := range c.RuntimeSection {
		errs = append(errs, ValidateRuntimePartial(c, r))
	}

	return errors2.JoinWithHead(errors.New("runtime section"), errs...)
}

func ValidateRuntimePartial(c Configuration, r Runtime) error {
	if err := r.Validate(c); err != nil {
		return errors2.JoinWithHead(fmt.Errorf("validating runtime '%s': %w", r.Name(), err), err)
	}
	return nil
}

func ValidateBuildSectionPartial(c Configuration) error {
	var errs []error

	for _, b := range c.BuildSection {
		errs = append(errs, ValidateBuildPartial(c, b))
	}

	return errors2.JoinWithHead(errors.New("build section"), errs...)
}

func ValidateBuildPartial(c Configuration, b Build) error {
	var errs []error

	foundBuildRuntime := false
	foundOuputRuntime := b.OutputRuntime == nil

	for _, r := range c.RuntimeSection {
		if r.Name() == b.BuildRuntime {
			foundBuildRuntime = true
		}
		if foundOuputRuntime && r.Name() == *b.OutputRuntime {
			foundOuputRuntime = true
		}
	}

	if !foundBuildRuntime {
		errs = append(errs, fmt.Errorf("build runtime '%s' not found", b.BuildRuntime))
	}
	if !foundOuputRuntime {
		errs = append(errs, fmt.Errorf("output runtime '%s' not found", *b.OutputRuntime))
	}

	errs = append(errs, ValidateCommandsPartial(c, b.Commands))

	if err := errors2.JoinWithHead(fmt.Errorf("validating build '%s'", b.Name), errs...); err != nil {
		return err
	}

	return nil
}

func ValidateDeploymentSectionPartial(c Configuration) error {
	var errs []error

	for _, d := range c.DeploymentSection {
		errs = append(errs, ValidateDeploymentPartial(c, d))
	}

	return errors2.JoinWithHead(errors.New("deployment section"), errs...)
}

func ValidateDeploymentPartial(c Configuration, d Deployment) error {
	var errs []error

	existing := make([]string, len(c.WorkflowSection))
	for i, w := range c.WorkflowSection {
		existing[i] = w.Name
	}

	processNotFound := func(w string) error {
		return fmt.Errorf("workflow '%s' not found", w)
	}

	errs = append(errs, slices.ProcessSubsetDifference(existing, d.Workflows, processNotFound)...)

	for _, cmd := range d.Commands {
		errs = append(errs, ValidateCommandPartial(c, cmd))
	}

	if err := errors2.JoinWithHead(fmt.Errorf("validating deployment '%s'", d.Name), errs...); err != nil {
		return err
	}

	return nil
}

func ValidateWorkflowSectionPartial(c Configuration) error {
	var errs []error

	for _, w := range c.WorkflowSection {
		errs = append(errs, ValidateWorkflowPartial(c, w))
	}

	return errors2.JoinWithHead(errors.New("workflow section"), errs...)
}

func ValidateWorkflowPartial(c Configuration, w Workflow) error {
	var errs []error

	existingR := make([]string, len(c.RuntimeSection))
	for i, r := range c.RuntimeSection {
		existingR[i] = r.Name()
	}
	processRuntimeNotFound := func(r string) error {
		return fmt.Errorf("runtime '%s' not found", r)
	}

	existingT := make([]string, len(c.TestSection))
	for i, t := range c.TestSection {
		existingT[i] = t.Name
	}

	for _, g := range w.Groups {
		// validate runtimes
		errs = append(errs, slices.ProcessSubsetDifference(existingR, g.Runtimes, processRuntimeNotFound)...)

		// validate tests
		errs = append(errs, slices.ProcessSubsetDifference(existingT, g.Tests, processRuntimeNotFound)...)
	}

	if err := errors2.JoinWithHead(fmt.Errorf("validating workflow '%s'", w.Name), errs...); err != nil {
		return err
	}

	return nil
}

func ValidateConfigSectionPartial(c Configuration) error {
	var errs []error

	return errors2.JoinWithHead(errors.New("config section"), errs...)
}
