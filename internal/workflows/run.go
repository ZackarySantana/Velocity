package workflows

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/config"
)

func findTest(c config.Config, workflow config.YAMLWorkflow, test string) (config.YAMLTest, error) {
	for t, yamlTest := range c.Tests {
		if t == test {
			return yamlTest, nil
		}
	}
	return config.YAMLTest{}, fmt.Errorf("test %s not found", test)
}

func RunWorkflow(c config.Config, workflow config.YAMLWorkflow) error {
	for image, testNames := range workflow.Tests {
		for _, testName := range testNames {
			test, err := findTest(c, workflow, string(testName))
			if err != nil {
				return err
			}
			fmt.Printf("Running test '%s' on image '%s'\n", test.Name, image)
		}
	}
	return nil
}
