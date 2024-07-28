package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type RoutineSection []Routine

func (r *RoutineSection) Validate() error {
	if r == nil {
		return nil
	}
	catcher := catcher.New()
	for _, routine := range *r {
		catcher.Catch(routine.Validate())
	}
	return catcher.Resolve()
}

type Routine struct {
	Name string   `yaml:"name"`
	Jobs []string `yaml:"jobs"`
}

func (r *Routine) Validate() error {
	oops := oops.With("routine", *r)
	if r.Name == "" {
		return oops.Errorf("name is required")
	}
	if len(r.Jobs) == 0 {
		return oops.Errorf("jobs are required")
	}
	return nil
}
