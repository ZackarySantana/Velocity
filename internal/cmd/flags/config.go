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
		Usage: "set the config file path",
		Value: "velocity.yml",
	}
)

func SetConfig(_ context.Context, cmd *cli.Command) error {
	filepath := cmd.String(ConfigFlag.Name)
	file, err := os.ReadFile(filepath)
	oops := oops.In("config_flag").With("config_name", filepath)
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

	cmd.Metadata[ConfigFlag.Name] = c

	Logger(cmd).Debug("Using config", "filepath", filepath)
	return nil
}

func Config(cmd *cli.Command) *config.Config {
	return cmd.Metadata[ConfigFlag.Name].(*config.Config)
}
