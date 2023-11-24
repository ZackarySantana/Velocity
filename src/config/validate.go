package config

func (c *Config) Validate() error {
	return nil
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
