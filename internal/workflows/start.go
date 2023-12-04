package workflows

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/config"
)

func StartWorkflow(c config.Config, workflow config.YAMLWorkflow) error {
	// Hit server endpoint to start workflow
	// Upload the config to the server, and get back a workflow ID

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
