package befores

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/workflows"
	"github.com/zackarysantana/velocity/src/config"
)

func Workflow(ctx *cli.Context) error {
	providedWorkflow := ctx.Args().First()
	c, err := GetConfig(ctx)
	if err != nil {
		return err
	}

	var w *config.YAMLWorkflow
	for title, workflow := range c.Workflows {
		if title == providedWorkflow {
			w = &workflow
			break
		}
	}

	if w == nil {
		if providedWorkflow == "" {
			fmt.Println("No workflow provided. Selecting from list.")
		} else {
			fmt.Printf("Workflow %s not found. Selecting from list.\n", providedWorkflow)
		}
		workflow, err := workflows.GetWorkflow(c, "Please select a workflow: ")
		if err != nil {
			return err
		}
		w = &workflow
	}

	ctx.App.Metadata["workflow"] = w
	return nil
}

func GetWorkflow(c *cli.Context) (*config.YAMLWorkflow, error) {
	workflow, ok := c.App.Metadata["workflow"].(*config.YAMLWorkflow)
	if !ok {
		return nil, cli.Exit("error getting workflow", 1)
	}
	return workflow, nil
}
