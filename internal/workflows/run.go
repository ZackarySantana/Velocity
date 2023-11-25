package workflows

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/config"
)

func RunWorkflow(c config.Config, workflow config.YAMLWorkflow) error {
	for image, testNames := range workflow.Tests {
		for _, testName := range testNames {
			test, err := c.GetTest(string(testName))
			if err != nil {
				return err
			}
			fmt.Printf("Running test '%s' on image '%s'\n", test.Name, image)
		}
	}
	return nil
}
