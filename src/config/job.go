package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
	"github.com/zackarysantana/velocity/src/config/id"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/job"
)

type JobSection []Job

func (j *JobSection) validateSyntax() error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for _, job := range *j {
		catcher.Catch(job.error().Wrap(job.validateSyntax()))
	}
	return catcher.Resolve()
}

func (j *JobSection) validateIntegrity(c *Config) error {
	if j == nil {
		return nil
	}
	catcher := catcher.New()
	for _, job := range *j {
		catcher.Catch(job.error().Wrap(job.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (j *JobSection) error() oops.OopsErrorBuilder {
	return oops.In("job_section")
}

func (j *JobSection) ToEntities(ic id.Creator, ec *entities.ConfigEntity) []*job.Job {
	jobs := make([]*job.Job, 0)
	for _, jb := range *j {
		jobs = append(jobs, jb.ToEntity(ic, ec))
	}
	return jobs
}

type Job struct {
	Name   string   `yaml:"name"`
	Tests  []string `yaml:"tests"`
	Images []string `yaml:"images"`
}

func (j *Job) validateSyntax() error {
	catcher := catcher.New()
	catcher.ErrorWhen(j.Name == "", "name is required")
	catcher.ErrorWhen(len(j.Tests) == 0, "tests are required")
	catcher.ErrorWhen(len(j.Images) == 0, "images are required")
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
		catcher.ErrorWhen(!found, "test %s not found", testName)
	}
	for _, imageName := range j.Images {
		found := false
		for _, image := range config.Images {
			if image.Name == imageName {
				found = true
				break
			}
		}
		catcher.ErrorWhen(!found, "image %s not found", imageName)
	}
	return catcher.Resolve()
}

func (j *Job) error() oops.OopsErrorBuilder {
	return oops.With("job_name", j.Name)
}

func (j *Job) ToEntity(ic id.Creator, ec *entities.ConfigEntity) *job.Job {
	tests := make([]string, len(j.Tests))
	for i, testName := range j.Tests {
		for _, test := range ec.Tests {
			if test.Name == testName {
				tests[i] = test.Id
				break
			}
		}
	}
	images := make([]string, len(j.Images))
	for i, imageName := range j.Images {
		for _, image := range ec.Images {
			if image.Name == imageName {
				images[i] = image.Id
				break
			}
		}
	}
	return &job.Job{
		Id:     ic(),
		Name:   j.Name,
		Tests:  tests,
		Images: images,
	}
}
