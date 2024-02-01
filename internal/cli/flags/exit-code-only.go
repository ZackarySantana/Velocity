package flags

import "github.com/urfave/cli/v2"

const ExitCodeOnlyFlagName = "exit-code-only"

type exitCode struct {
	boolFlag
}

func ExitCodeOnly() exitCode {
	return exitCode{
		boolFlag: boolFlag{
			BoolFlag: cli.BoolFlag{
				Name:    ExitCodeOnlyFlagName,
				Aliases: []string{"eco"},
				Usage:   "only return the exit code",
			},
		},
	}
}
