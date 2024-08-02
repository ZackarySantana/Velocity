package flags

import (
	"context"
	"os"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/src/config"
)

var (
	ConfigFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "set the config file location",
		Value: "velocity.yml",
	}
)

func SetConfig(_ context.Context, cmd *cli.Command) error {
	file, err := os.ReadFile(cmd.String("config"))
	oops := oops.In("config_flag").With("config_name", cmd.String("config"))
	if err != nil {
		return oops.Wrap(err)
	}

	c, err := config.Parse(file)
	if err != nil {
		return oops.Wrap(err)
	}

	if err := c.Validate(); err != nil {
		return oops.Wrap(err)
	}

	cmd.Metadata["config"] = c

	Logger(cmd).Debug("Using config", "location", cmd.String("config"))
	return nil
}

func Config(cmd *cli.Command) *config.Config {
	return cmd.Metadata["config"].(*config.Config)
}
