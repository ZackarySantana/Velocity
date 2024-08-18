package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Tests    TestSection    `yaml:"tests"`
	Images   ImageSection   `yaml:"images"`
	Jobs     JobSection     `yaml:"jobs"`
	Routines RoutineSection `yaml:"routines"`
}

func (c *Config) Validate() error {
	if c == nil {
		return oops.Errorf("config is nil")
	}
	catcher := catcher.New()
	catcher.Catch(Validate(&c.Tests, c))
	catcher.Catch(Validate(&c.Images, c))
	catcher.Catch(Validate(&c.Jobs, c))
	catcher.Catch(Validate(&c.Routines, c))
	return catcher.Resolve()
}

func (c *Config) GetJob(jobName string) (Job, error) {
	for _, job := range c.Jobs {
		if job.Name == jobName {
			return job, nil
		}
	}
	return Job{}, oops.Errorf("job not found: %s", jobName)
}

func (c *Config) GetImage(imageName string) (Image, error) {
	for _, image := range c.Images {
		if image.Name == imageName {
			return image, nil
		}
	}
	return Image{}, oops.Errorf("image not found: %s", imageName)
}

func (c *Config) GetTest(testName string) (Test, error) {
	for _, test := range c.Tests {
		if test.Name == testName {
			return test, nil
		}
	}
	return Test{}, oops.Errorf("test not found: %s", testName)
}

func Parse(bytes []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}
