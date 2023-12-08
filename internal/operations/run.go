package operations

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/befores"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/internal/workflows"
)

var Run = []*cli.Command{
	{
		Name:      "run",
		Aliases:   []string{"r"},
		Usage:     "run a workflow",
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

			fmt.Println("Running workflow " + w.Name)

			return workflows.StartWorkflow(c, *w)
		},
	},
}
