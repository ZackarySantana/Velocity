package flags

import (
	"github.com/urfave/cli/v2"
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
