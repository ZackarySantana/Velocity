package flags

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	configPackage "github.com/zackarysantana/velocity/src/config"
	"github.com/zackarysantana/velocity/src/config/configuration"
)

const ConfigFlagName = "config"

type config struct {
	stringFlag
}

func Config() config {
	return config{
		stringFlag: stringFlag{
			StringFlag: cli.StringFlag{
				Name:    ConfigFlagName,
				Aliases: []string{"c"},
				Usage:   "location of your configuration file",
				Value:   "velocity.yml",
			},
		},
	}
}

func ParseConfigFromFlag(ctx *cli.Context) (*configuration.Configuration, error) {
	fileName := ctx.String(ConfigFlagName)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading your file '%s': %v", fileName, err)
	}

	raw, err := configPackage.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("there was an error parsing your file in to yaml: %v", err)
	}

	parsed, err := configPackage.HydrateConfiguration(raw)
	if err != nil {
		return nil, fmt.Errorf("there was an error hydrating your config: %v", err)
	}

	err = configPackage.ValidateConfiguration(*parsed)
	if err != nil {
		return nil, fmt.Errorf("your configuration is invalid, see:\n %v", err)
	}

	return parsed, nil
}
