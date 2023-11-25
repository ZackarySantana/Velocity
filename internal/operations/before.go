package operations

import (
	"os"

	"github.com/urfave/cli"
	"github.com/zackarysantana/velocity/src/config"
)

func BeforeConfig(c *cli.Context) error {
	if c.String("config") != "" {
		os.Setenv("VELOCITY_CONFIG", c.String("config"))
	}
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.App.Metadata["config"] = *config
	return err
}

func GetConfig(c *cli.Context) (config.Config, error) {
	config, ok := c.App.Metadata["config"].(config.Config)
	if !ok {
		return config, cli.NewExitError("config not found", 1)
	}
	return config, nil
}

func CombineBefores(beforeFuncs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, beforeFunc := range beforeFuncs {
			err := beforeFunc(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
