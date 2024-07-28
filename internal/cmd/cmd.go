package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:                  "velocity",
		Version:               "0.0.1",
		Description:           "CLI for your CI",
		EnableShellCompletion: true,
		Metadata:              map[string]interface{}{},
		Flags: []cli.Flag{
			loggerModeFlag,
			verboseFlag,
		},
		Before: func(ctx context.Context, cmd *cli.Command) error {
			if err := setLogger(ctx, cmd); err != nil {
				return err
			}
			return nil
		},
		Commands: []*cli.Command{
			routine,
		},
		ExitErrHandler: func(ctx context.Context, cmd *cli.Command, err error) {
			oops, _ := oops.AsOops(err)
			fmt.Println(oops)
			Logger(cmd).Error(
				"Exiting with error",
				slog.Any("error", oops.LogValuer()),
			)
		},
	}
}
