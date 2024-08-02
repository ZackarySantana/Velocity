package cmd

import (
	"context"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/internal/cmd/flags"
)

var (
	routine = &cli.Command{
		Name:  "routine",
		Usage: "routine related commands",
		Commands: []*cli.Command{
			routineRun,
		},
	}
	routineRun = &cli.Command{
		Name:  "run",
		Usage: "runs a routine locally",
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
		Before: befores(flags.SetConfig),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := flags.Config(cmd)

			flags.Logger(cmd).Info("Tests", "tests", c.Tests)

			return oops.Code("Testing").Errorf("Not implemented")
		},
	}
)
