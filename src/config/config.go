package config

import (
	"fmt"
)

func (c *Config) GetWorkflow(w string) (YAMLWorkflow, error) {
	for name, workflow := range c.Workflows {
		if name == w {
			return workflow, nil
		}
	}
	return YAMLWorkflow{}, fmt.Errorf("workflow %s not found", w)
}

func (c *Config) GetImage(i string) (YAMLImage, error) {
	for name, image := range c.Images {
		if name == i {
			return image, nil
		}
	}
	return YAMLImage{}, fmt.Errorf("image %s not found", i)
}

func (c *Config) GetTest(t string) (YAMLTest, error) {
	for name, test := range c.Tests {
		if name == t {
			return test, nil
		}
	}
	return YAMLTest{}, fmt.Errorf("test %s not found", t)
}
