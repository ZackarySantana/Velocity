package operations

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zackarysantana/velocity/internal/workflows"
	"github.com/zackarysantana/velocity/src/config"
)

var Run = []cli.Command{
	{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "run a workflow",
		ArgsUsage: "[workflow]",
		Flags: []cli.Flag{
			configFlag,
		},
		Before: CombineBefores(BeforeConfig),
		Action: func(cCtx *cli.Context) error {
			providedWorkflow := cCtx.Args().First()
			c, err := GetConfig(cCtx)
			if err != nil {
				return err
			}

			var w *config.YAMLWorkflow
			for title, workflow := range c.Workflows {
				if title == providedWorkflow {
					w = &workflow
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

			fmt.Println("Running workflow " + w.Name)

			err = workflows.StartWorkflow(c, *w)
			if err != nil {
				return err
			}

			return nil
		},
	},
}
