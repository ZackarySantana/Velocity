package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
	"github.com/zackarysantana/velocity/src/config/id"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/routine"
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

func (r *RoutineSection) ToEntities(ic id.Creator, ec *entities.ConfigEntity) []*routine.Routine {
	routines := make([]*routine.Routine, 0)
	for _, rt := range *r {
		routines = append(routines, rt.ToEntity(ic, ec))
	}
	return routines
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

func (r *Routine) ToEntity(ic id.Creator, ec *entities.ConfigEntity) *routine.Routine {
	jobs := make([]string, len(r.Jobs))
	for i, jobName := range r.Jobs {
		for _, job := range ec.Jobs {
			if job.Name == jobName {
				jobs[i] = job.Id
				break
			}
		}
	}
	return &routine.Routine{
		Id:   ic(),
		Name: r.Name,
		Jobs: jobs,
	}
}
