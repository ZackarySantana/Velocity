package operations

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/internal/operations/befores"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/src/clients/v1types"
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
		Before: befores.CombineBefores(befores.Config, befores.Workflow, befores.VelocityClient),
		Action: func(ctx *cli.Context) error {
			c, err := befores.GetConfig(ctx)
			if err != nil {
				return err
			}

			w, err := befores.GetWorkflow(ctx)
			if err != nil {
				return err
			}

			client, err := befores.GetVelocityClient(ctx)
			if err != nil {
				return err
			}

			gitCtx, err := jobs.NewCurrentContext()
			if err != nil {
				return err
			}

			fmt.Println("Running workflow " + w.Name)
			req := v1types.PostInstanceStartRequest{
				ProjectName: &c.Config.Project,
				Config:      c,
				Workflow:    w.Name,
				GitHash:     gitCtx.CommitHash,
			}
			resp, err := client.PostInstanceStart(req)
			if err != nil {
				return err
			}

			fmt.Println("Instance ID: " + resp.InstanceId)

			if c.Config.UI != nil {
				fmt.Printf("Open instance at %s/instances/%s\n", *c.Config.UI, resp.InstanceId)
			}

			return nil
		},
	},
}
