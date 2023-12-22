package operations

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type CLIApp struct {
	app *cli.App
}

func NewCLIApp() CLIApp {
	cliApp := CLIApp{}
	cliApp.app = &cli.App{
		Name:     "velocity",
		Version:  "0.0.1",
		Usage:    "manage, run, and report on tests quickly",
		Commands: appendCommands(Run, Validate, RunLocal, FirstTimeRegister),
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			fmt.Println(err)
		},
	}
	return cliApp
}

func (c CLIApp) Run() error {
	return c.app.Run(os.Args)
}

func appendCommands(commands ...[]*cli.Command) []*cli.Command {
	var allCommands []*cli.Command
	for _, command := range commands {
		allCommands = append(allCommands, command...)
	}
	return allCommands
}
