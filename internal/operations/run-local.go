package operations

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/befores"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/internal/workflows"
)

var RunLocal = []*cli.Command{
	{
		Name:      "run-local",
		Aliases:   []string{"rl"},
		Usage:     "run a workflow locally",
		ArgsUsage: "[workflow]",
		Flags: []cli.Flag{
			flags.Config,
		},
		Before: befores.CombineBefores(befores.Config, befores.Workflow),
		Action: func(ctx *cli.Context) error {
			c, err := befores.GetConfig(ctx)
			if err != nil {
				return err
			}

			w, err := befores.GetWorkflow(ctx)
			if err != nil {
				return err
			}

			fmt.Println("Running workflow synchronously" + w.Name)

			results, err := workflows.RunSyncWorkflow(c, *w)
			if err != nil {
				return err
			}

			fmt.Println("\n\nWorkflow completed: ")

			for _, result := range results {
				if result.Success != nil {
					fmt.Print("✅")
					fmt.Println("'" + result.Job.GetImage() + "' ran '" + (result.Job.GetCommand()) + "'")
					fmt.Println(result.Success.Logs)
					fmt.Println()
				}

				if result.Failed != nil {
					fmt.Println("❌")
					fmt.Println("'" + result.Job.GetImage() + "' ran '" + (result.Job.GetCommand()) + "'")
					fmt.Println(result.Failed.Error)
					fmt.Println()
				}
			}

			return nil
		},
	},
}
