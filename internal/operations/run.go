package operations

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/befores"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/internal/workflows"
	"github.com/zackarysantana/velocity/src/config"
)

var Run = []*cli.Command{
	{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "run a workflow",
		ArgsUsage: "[workflow]",
		Flags: []cli.Flag{
			flags.Config,
			flags.Sync,
		},
		Before: befores.CombineBefores(befores.Config, befores.Sync),
		Action: func(cCtx *cli.Context) error {
			providedWorkflow := cCtx.Args().First()
			c, err := befores.GetConfig(cCtx)
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

			sync := befores.GetSync(cCtx)

			if sync {
				fmt.Println("Running workflow in sync")
				results, err := workflows.RunSyncWorkflow(c, *w)

				if err != nil {
					return err
				}

				fmt.Print("\n\nWorkflow completed: ")

				for _, result := range results {
					if result.Success != nil {
						fmt.Println("'" + result.Job.Image + "' ran '" + result.Job.Command + "'")
						fmt.Println(result.Success.Logs)
						fmt.Println()
					}

					if result.Failed != nil {
						fmt.Println("'" + result.Job.Image + "' ran '" + result.Job.Command + "'")
						fmt.Println(result.Failed.Error)
						fmt.Println()
					}
				}

				return nil

			}
			return workflows.StartWorkflow(c, *w)
		},
	},
}
