package cmd

import (
	"context"
	"fmt"
	"log/slog"
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
			routineList,
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

			flags.Logger(cmd).Info("routine",
				slog.String("id", fmt.Sprintf("%v", data.Id)),
				slog.String("link", flags.API(cmd).GetLink().Routine(ctx, fmt.Sprintf("%s", data.Id))),
				slog.String("routine", routine),
			)

			return nil
		},
	}
	routineList = &cli.Command{
		Name:  "list",
		Usage: "lists all routines",
		Flags: []cli.Flag{
			flags.ConfigFlag,
		},
		Before: befores(flags.SetConfig),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := flags.Config(cmd)

			routinesAttrs := []slog.Attr{}

			for _, routine := range c.Routines {
				jobAttrs := []any{}

				for _, jobName := range routine.Jobs {
					job, err := c.GetJob(jobName)
					if err != nil {
						return oops.Code("routine").With("job_name", jobName).Wrap(err)
					}

					testAttrs := []any{}
					for _, testName := range job.Tests {
						testAttrs = append(testAttrs, slog.String("name", testName))
					}

					imageAttrs := []any{}
					for _, imageName := range job.Images {
						imageAttrs = append(imageAttrs, slog.String("name", imageName))
					}

					jobAttrs = append(jobAttrs, slog.Group(jobName,
						slog.Group("tests", testAttrs...),
						slog.Group("images", imageAttrs...),
					))
				}

				routineAttr := slog.Group(routine.Name, slog.Group("jobs", jobAttrs...))
				routinesAttrs = append(routinesAttrs, routineAttr)
			}

			flags.Logger(cmd).LogAttrs(ctx, slog.LevelInfo, "routines", routinesAttrs...)

			return nil
		},
	}
)
