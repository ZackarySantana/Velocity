package befores

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/src/config"
)

func Config(c *cli.Context) error {
	if c.String(flags.Config.Name) != "" {
		os.Setenv("VELOCITY_CONFIG", c.String(flags.Config.Name))
	}
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.App.Metadata[flags.Config.Name] = config
	return err
}

func GetConfig(c *cli.Context) (*config.Config, error) {
	config, ok := c.App.Metadata[flags.Config.Name].(*config.Config)
	if !ok {
		return nil, cli.Exit("config not found", 1)
	}
	return config, nil
}
