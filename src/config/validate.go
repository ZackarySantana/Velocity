package config

import (
	"errors"
	"fmt"

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

	return errors.Join(errs...)
}

func ValidateTestSectionPartial(c Configuration) error {
	var errs []error

	for _, t := range c.TestSection {
		errs = append(errs, ValidateTestPartial(c, t))
	}

	return errors.Join(errs...)
}

func ValidateTestPartial(c Configuration, t Test) error {
	var errs []error

	for _, cmd := range t.Commands {
		errs = append(errs, ValidateCommandPartial(c, cmd))
	}

	return errors.Join(errs...)
}

func ValidateCommandPartial(c Configuration, cmd Command) error {
	return cmd.Validate(c)
}

func ValidateOperationSectionPartial(c Configuration) error {
	var errs []error

	return errors.Join(errs...)
}

func ValidateOperationPartial(c Configuration, o Operation) error {
	var errs []error

	for _, cmd := range o.Commands {
		errs = append(errs, ValidateCommandPartial(c, cmd))
	}

	return errors.Join(errs...)
}

func ValidateRuntimeSectionPartial(c Configuration) error {
	var errs []error

	for _, r := range c.RuntimeSection {
		errs = append(errs, ValidateRuntimePartial(c, r))
	}

	return errors.Join(errs...)
}

func ValidateRuntimePartial(c Configuration, r Runtime) error {
	return r.Validate(c)
}

func ValidateBuildSectionPartial(c Configuration) error {
	var errs []error

	for _, b := range c.BuildSection {
		errs = append(errs, ValidateBuildPartial(c, b))
	}

	return errors.Join(errs...)
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

	for _, cmd := range b.Commands {
		errs = append(errs, ValidateCommandPartial(c, cmd))
	}

	return errors.Join(errs...)
}

func ValidateDeploymentSectionPartial(c Configuration) error {
	var errs []error

	for _, d := range c.DeploymentSection {
		errs = append(errs, ValidateDeploymentPartial(c, d))
	}

	return errors.Join(errs...)
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

	return errors.Join(errs...)
}

func ValidateWorkflowSectionPartial(c Configuration) error {
	var errs []error

	for _, w := range c.WorkflowSection {
		errs = append(errs, ValidateWorkflowPartial(c, w))
	}

	return errors.Join(errs...)
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

	return errors.Join(errs...)
}

func ValidateConfigSectionPartial(c Configuration) error {
	var errs []error

	return errors.Join(errs...)
}

type Commands []Command

func IsPrebuiltCommand(name string) bool {
	fmt.Println("IsPrebuiltCommand", name)
	return true
}

func (o Operation) Validate(c Configuration) error {
	return ValidateCommands(c, o.Commands)
}

func (d *Deployment) Validate(c Configuration) error {
	var errs []error

	// Validate workflows inside of deployment
	for _, w := range d.Workflows {
		found := false
		for _, w2 := range c.WorkflowSection {
			if w == w2.Name {
				found = true
				break
			}
		}
		if !found {
			errs = append(errs, fmt.Errorf("workflow '%s' not found", w))
		}
	}

	// Validate commands inside of deployment
	errs = append(errs, ValidateCommands(c, d.Commands))

	return errors.Join(errs...)
}

func ValidateCommands(c Configuration, cmds []Command) error {
	var errs []error
	for _, w := range cmds {
		errs = append(errs, w.Validate(c))
	}
	return errors.Join(errs...)
}
