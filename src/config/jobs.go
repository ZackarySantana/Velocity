package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type JobSection []Job

func (j *JobSection) Validate() error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for _, job := range *j {
		catcher.Catch(job.Validate())
	}
	return catcher.Resolve()
}

type Job struct {
	Name   string   `yaml:"name"`
	Tests  []string `yaml:"tests"`
	Images []string `yaml:"images"`
}

func (j *Job) Validate() error {
	oops := oops.With("job", *j)
	if j.Name == "" {
		return oops.Errorf("name is required")
	}
	if len(j.Tests) == 0 {
		return oops.Errorf("tests are required")
	}
	if len(j.Images) == 0 {
		return oops.Errorf("images are required")
	}
	return nil
}
