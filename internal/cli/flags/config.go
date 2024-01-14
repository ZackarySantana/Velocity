package flags

import (
	"github.com/urfave/cli/v2"
)

type config struct {
	sf cli.StringFlag
}

func Config() config {
	cf := config{
		sf: cli.StringFlag{
			Name:      "config",
			Aliases:   []string{"c"},
			Usage:     "location of your configuration file",
			TakesFile: true,
			Value:     "velocity.yml",
		},
	}
	return cf
}

func (cf config) Flag() *cli.StringFlag {
	return &cf.sf
}
