package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/internal/stats"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:                  "velocity",
		Version:               "0.0.1",
		Description:           "CLI for your CI",
		EnableShellCompletion: true,
		Metadata:              map[string]interface{}{},
		Before: func(ctx context.Context, cmd *cli.Command) error {
			// TODO: Custom handlers for different providers
			// one that is just plain text for each message.
			opts := slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}
			logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
			cmd.Metadata["logger"] = logger

			logger.Debug("Starting CLI", "cmd", os.Args, "version", cmd.Version, "ip", stats.GetIP())
			return nil
		},
		Commands: []*cli.Command{
			routine,
		},
	}
}

func Logger(cmd *cli.Command) *slog.Logger {
	return cmd.Root().Metadata["logger"].(*slog.Logger)
}
