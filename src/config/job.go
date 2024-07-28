package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type JobSection []Job

func (j *JobSection) validateSyntax() error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for _, job := range *j {
		catcher.Catch(job.Error().Wrap(job.validateSyntax()))
	}
	return catcher.Resolve()
}

func (j *JobSection) validateIntegrity(c *Config) error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for _, job := range *j {
		catcher.Catch(job.Error().Wrap(job.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (j *JobSection) Error() oops.OopsErrorBuilder {
	return oops.Code("job_section")
}

type Job struct {
	Name   string   `yaml:"name"`
	Tests  []string `yaml:"tests"`
	Images []string `yaml:"images"`
}

func (j *Job) validateSyntax() error {
	catcher := catcher.New()
	if j.Name == "" {
		catcher.Error("name is required")
	}
	if len(j.Tests) == 0 {
		catcher.Error("tests are required")
	}
	if len(j.Images) == 0 {
		catcher.Error("images are required")
	}
	return catcher.Resolve()
}

func (j *Job) validateIntegrity(config *Config) error {
	catcher := catcher.New()
	for _, testName := range j.Tests {
		found := false
		for _, test := range config.Tests {
			if test.Name == testName {
				found = true
				break
			}
		}
		if !found {
			catcher.Error("test %s not found", testName)
		}
	}
	for _, imageName := range j.Images {
		found := false
		for _, image := range config.Images {
			if image.Name == imageName {
				found = true
				break
			}
		}
		if !found {
			catcher.Error("image %s not found", imageName)
		}
	}
	return catcher.Resolve()
}

func (j *Job) Error() oops.OopsErrorBuilder {
	return oops.With("job_name", j.Name)
}
