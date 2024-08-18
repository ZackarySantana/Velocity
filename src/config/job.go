package config

import (
	"fmt"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type JobSection []Job

func (j *JobSection) validateSyntax() error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for i, job := range *j {
		catcher.Catch(job.error(i).Wrap(job.validateSyntax()))
	}
	return catcher.Resolve()
}

func (j *JobSection) validateIntegrity(c *Config) error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for i, job := range *j {
		catcher.Catch(job.error(i).Wrap(job.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (j *JobSection) error(_ int) oops.OopsErrorBuilder {
	return oops.In("job_section")
}

type Job struct {
	Name   string   `yaml:"name"`
	Tests  []string `yaml:"tests"`
	Images []string `yaml:"images"`
}

func (j *Job) validateSyntax() error {
	catcher := catcher.New()
	catcher.When(j.Name == "", "name is required")
	catcher.When(len(j.Tests) == 0, "tests are required")
	catcher.When(len(j.Images) == 0, "images are required")
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
		catcher.When(!found, "test %s not found", testName)
	}
	for _, imageName := range j.Images {
		found := false
		for _, image := range config.Images {
			if image.Name == imageName {
				found = true
				break
			}
		}
		catcher.When(!found, "image %s not found", imageName)
	}
	return catcher.Resolve()
}

func (j *Job) error(i int) oops.OopsErrorBuilder {
	return oops.With(fmt.Sprintf("job_name_%d", i), j.Name)
}
