package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/cli/flags"
)

func CreateWorkflowList(app *app) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list all workflows",
		Flags: []cli.Flag{
			flags.Config().Flag(),
		},
		Action: func(ctx *cli.Context) error {
			config, err := flags.ParseConfigFromFlag(ctx)
			if err != nil {
				return err
			}

			for _, w := range config.WorkflowSection {
				app.logf(ctx, "Workflow: %s", w.Name)
				for _, g := range w.Groups {
					app.logf(ctx, "    Group: %s", g.Name)
					app.logf(ctx, "        Images: %s", g.Runtimes)
					app.logf(ctx, "        Tests: %s", g.Tests)
				}
			}

			return nil
		},
	}
}

func CreateWorkflow(app *app) *cli.Command {
	return &cli.Command{
		Name:    "workflow",
		Aliases: []string{"w"},
		Usage:   "manage, run, and view your workflows",
		Subcommands: []*cli.Command{
			CreateWorkflowList(app),
		},
	}
}
