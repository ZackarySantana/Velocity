package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/cli/flags"
)

func CreateValidate(app app) *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"v"},
		Usage:   "validate a configuration file",
		Flags: []cli.Flag{
			flags.Config().Flag(),
			flags.ExitCodeOnly().Flag(),
		},
		Action: func(ctx *cli.Context) error {
			if _, err := flags.ParseConfigFromFlag(ctx); err != nil {
				return err
			}

			fmt.Printf("Configuration (%s) is valid!\n", ctx.String("config"))

			return nil
		},
	}
}
