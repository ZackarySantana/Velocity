package cmd

import (
	"context"
	"log/slog"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/internal/cmd/flags"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:                  "velocity",
		Version:               "0.0.1",
		Description:           "CLI for your CI",
		EnableShellCompletion: true,
		Metadata:              map[string]interface{}{},
		Flags: []cli.Flag{
			flags.LoggerModeFlag,
			flags.VerboseFlag,
		},
		Before: befores(flags.SetLogger),
		Commands: []*cli.Command{
			routine,
		},
		ExitErrHandler: func(ctx context.Context, cmd *cli.Command, err error) {
			oops, ok := oops.AsOops(err)
			if ok {
				flags.Logger(cmd).Error(
					"Exiting with error",
					slog.Any("error", oops.LogValuer()),
				)
				return
			}
			flags.Logger(cmd).Error("Exiting with error", "error", err)
		},
	}
}

func befores(fns ...cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		for _, fn := range fns {
			if err := fn(ctx, cmd); err != nil {
				return err
			}
		}
		return nil
	}
}
