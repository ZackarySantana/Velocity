package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type RoutineSection []Routine

func (r *RoutineSection) validateSyntax() error {
	if r == nil {
		return nil
	}
	catcher := catcher.New()
	for _, routine := range *r {
		catcher.Catch(routine.error().Wrap(routine.validateSyntax()))
	}
	return catcher.Resolve()
}

func (r *RoutineSection) validateIntegrity(c *Config) error {
	if r == nil {
		return nil
	}
	catcher := catcher.New()
	for _, routine := range *r {
		catcher.Catch(routine.error().Wrap(routine.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (r *RoutineSection) error() oops.OopsErrorBuilder {
	return oops.In("routine_section")
}

func (r *RoutineSection) GetRoutine(name string) *Routine {
	for _, routine := range *r {
		if routine.Name == name {
			return &routine
		}
	}
	return nil
}

type Routine struct {
	Name string   `yaml:"name"`
	Jobs []string `yaml:"jobs"`
}

func (r *Routine) validateSyntax() error {
	catcher := catcher.New()
	if r.Name == "" {
		catcher.Error("name is required")
	}
	if len(r.Jobs) == 0 {
		catcher.Error("jobs are required")
	}
	return catcher.Resolve()
}

func (r *Routine) validateIntegrity(config *Config) error {
	catcher := catcher.New()
	for _, jobName := range r.Jobs {
		found := false
		for _, job := range config.Jobs {
			if job.Name == jobName {
				found = true
				break
			}
		}
		catcher.ErrorWhen(!found, "job %s not found", jobName)
	}
	return catcher.Resolve()
}

func (r *Routine) error() oops.OopsErrorBuilder {
	return oops.With("routine_name", r.Name)
}
