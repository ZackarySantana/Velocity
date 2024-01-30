package flags

import (
	"errors"

	"github.com/urfave/cli/v2"
)

const VerboseFlagName = "verbose"

type verbose struct {
	stringFlag
}

func Verbose() verbose {
	return verbose{
		stringFlag: stringFlag{
			StringFlag: cli.StringFlag{
				Name:    VerboseFlagName,
				Aliases: []string{"v"},
				Usage:   "lowest level of verbosity for the output (info, warning, or error)",
				Value:   "error",
				Action: func(ctx *cli.Context, v string) error {
					if v == "info" || v == "warning" || v == "error" {
						return nil
					}
					return errors.New("verbosity can only be info, warning, or error")
				},
			},
		},
	}
}
