package config

import (
	"errors"
	"fmt"
)

func (c *Config) validateTest(t YAMLTest) error {
	if t.Language == nil && t.Framework == nil {
		if t.Run == nil {
			return fmt.Errorf("test '%s' must have either language & framework or run", t.Name)
		}
		return nil
	}

	if t.Language != nil && t.Framework != nil {
		if t.Run != nil {
			return fmt.Errorf("test '%s' must have either language & framework or run- not both", t.Name)
		}
		return nil
	}

	if t.Language != nil || t.Framework != nil {
		return fmt.Errorf("test '%s' must have both language & framework or run- not both", t.Name)
	}

	if t.Run == nil {
		return fmt.Errorf("test '%s' must have either language & framework or run- not both", t.Name)
	}

	return nil
}

func (c *Config) validateImage(i YAMLImage) error {
	return nil
}

func (c *Config) validateWorkflow(w YAMLWorkflow) error {
	errs := []error{}
	for image, tests := range w.Tests {
		_, err := c.GetImage(image)
		if err != nil {
			errs = append(errs, fmt.Errorf("image '%s' is not defined in config but used in '%s' workflow", image, w.Name))
		}
		for _, test := range tests {
			_, err = c.GetTest(string(test))
			if err != nil {
				errs = append(errs, fmt.Errorf("test '%s' is not defined in config but used in '%s' workflow inside image '%s'", test, w.Name, image))
			}
		}
	}
	return errors.Join(errs...)
}

func (c *Config) Validate() error {
	errs := []error{}
	for _, test := range c.Tests {
		if err := c.validateTest(test); err != nil {
			errs = append(errs, err)
		}
	}
	for _, image := range c.Images {
		if err := c.validateImage(image); err != nil {
			errs = append(errs, err)
		}
	}
	for _, workflow := range c.Workflows {
		if err := c.validateWorkflow(workflow); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (c *Config) Populate() error {
	for title, workflow := range c.Workflows {
		workflow.Name = title
		c.Workflows[title] = workflow
	}
	for title, test := range c.Tests {
		test.Name = title
		c.Tests[title] = test
	}
	for title, image := range c.Images {
		image.Name = title
		c.Images[title] = image
	}
	return nil
}
