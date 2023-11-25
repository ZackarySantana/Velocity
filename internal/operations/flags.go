package operations

import "github.com/urfave/cli"

var (
	configFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "location of your configuration file",
	}
)
