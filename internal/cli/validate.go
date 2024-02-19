package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/cli/flags"
)

func CreateValidate(app *app) *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"v"},
		Usage:   "validate a configuration file",
		Flags: []cli.Flag{
			flags.Config().Flag(),
		},
		Action: func(ctx *cli.Context) error {
			if _, err := flags.ParseConfigFromFlag(ctx); err != nil {
				return err
			}

			app.logf(ctx, "Configuration (%s) is valid!", ctx.String("config"))

			return nil
		},
	}
}
