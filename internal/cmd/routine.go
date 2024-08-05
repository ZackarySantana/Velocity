package cmd

import (
	"context"
	"strings"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/internal/cmd/flags"
	"github.com/zackarysantana/velocity/src/velocity"
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
			flags.ServerFlag,
		},
		Before: befores(flags.SetConfig, flags.SetServer),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := flags.Config(cmd)
			routine := strings.Join(cmd.Args().Slice(), " ")
			_, data, err := velocity.New(flags.Server(cmd)).StartRoutine(c, routine)
			if err != nil {
				return oops.Code("request").Wrap(err)
			}

			flags.Logger(cmd).Info(data.Id)

			return nil
		},
	}
)
