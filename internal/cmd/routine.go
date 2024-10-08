package cmd

import (
	"context"
	"fmt"
	"strings"

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
			flags.APIFlag,
		},
		Before: befores(flags.SetConfig, flags.SetAPI),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := flags.Config(cmd)
			routine := strings.Join(cmd.Args().Slice(), " ")
			_, data, err := flags.API(cmd).StartRoutine(ctx, c, routine)
			if err != nil {
				return oops.Code("request").Wrap(err)
			}

			flags.Logger(cmd).Info(fmt.Sprintf("%v", data.Id))

			return nil
		},
	}
)
