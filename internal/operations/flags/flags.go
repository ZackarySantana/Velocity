package flags

import "github.com/urfave/cli/v2"

var (
	Config = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "location of your configuration file",
	}
	Sync = &cli.BoolFlag{
		Name:  "sync",
		Usage: "sync your configuration file with the remote",
	}
)
