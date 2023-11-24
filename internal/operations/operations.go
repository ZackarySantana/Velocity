package operations

import (
	"os"

	"github.com/urfave/cli"
)

type CLIApp struct {
	app *cli.App
}

func NewCLIApp() CLIApp {
	cliApp := CLIApp{}
	cliApp.app = &cli.App{
		Name:     "velocity",
		Usage:    "manage, run, and report on tests quickly",
		Commands: append([]cli.Command{}, Run...),
	}
	return cliApp
}

func (c CLIApp) Run() error {
	return c.app.Run(os.Args)
}
