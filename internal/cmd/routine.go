package cmd

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/src/config"
)

var (
	routine = &cli.Command{
		Name:  "routine",
		Usage: "add a task to the list",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			file, err := os.ReadFile("velocity.yml")
			if err != nil {
				return err
			}

			c, err := config.Read(file)
			if err != nil {
				return err
			}

			Logger(cmd).Debug("Adding routine")
			Logger(cmd).Info("Tests", "tests", c.Tests)

			return nil
		},
	}
)
