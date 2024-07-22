package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/samber/oops"
	slogmulti "github.com/samber/slog-multi"
	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/internal/stats"
	"github.com/zackarysantana/velocity/internal/vlog"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:                  "velocity",
		Version:               "0.0.1",
		Description:           "CLI for your CI",
		EnableShellCompletion: true,
		Metadata:              map[string]interface{}{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "mode",
				Usage: "set the mode to `debug|quiet`",
				Validator: func(s string) error {
					if s == "" {
						return nil
					}
					if s != "debug" && s != "quiet" {
						return oops.In("flags").Errorf("mode must be debug or quiet")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose output",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) error {
			// TODO: Custom handlers for different providers
			// one that is just plain text for each message.
			level := slog.LevelInfo
			switch cmd.String("mode") {
			case "debug":
				level = slog.LevelDebug
			case "error":
				level = slog.LevelError
			case "quiet":
				level = 12
			}
			var stdLogger slog.Handler
			if cmd.Bool("verbose") {
				stdLogger = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
			} else {
				stdLogger = vlog.NewPlainHandler(os.Stdout, &vlog.Options{Level: level})
			}

			logger := slog.New(slogmulti.Fanout(
				stdLogger,
			))
			cmd.Metadata["logger"] = logger

			logger.Debug("Starting logger", "cmd", os.Args, "version", cmd.Version, "ip", stats.GetIP())
			return nil
		},
		Commands: []*cli.Command{
			routine,
		},
		ExitErrHandler: func(ctx context.Context, cmd *cli.Command, err error) {
			Logger(cmd).Error(
				"Exiting with error",
				slog.Any("error", err), // unwraps and flattens error context
			)
		},
	}
}

func Logger(cmd *cli.Command) *slog.Logger {
	return cmd.Root().Metadata["logger"].(*slog.Logger)
}
