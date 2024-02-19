package flags

import (
	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/gen/pkl/velocity"
	config_ "github.com/zackarysantana/velocity/src/config"
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
				Value:   "velocity.pkl",
			},
		},
	}
}

func ParseConfigFromFlag(ctx *cli.Context) (*velocity.Velocity, error) {
	return config_.Load(ctx.Context, ctx.String(ConfigFlagName))
}
