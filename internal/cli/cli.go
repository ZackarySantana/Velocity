package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/cli/flags"
)

type app struct {
	exitOnError bool
	silent      bool
}

func createCliApp(app *app) *cli.App {
	return &cli.App{
		Name:    "velocity",
		Version: "0.0.1",
		Usage:   "manage, run, and report on tests quickly",
		Flags: []cli.Flag{
			flags.ExitCodeOnly().Flag(),
		},
		Commands: []*cli.Command{
			CreateValidate(app),
			CreateWorkflow(app),
		},
		ExitErrHandler: app.exitErrHandler,
	}
}

func Run() error {
	app := app{}
	return createCliApp(&app).Run(os.Args)
}

func RunAndExitOnError() {
	app := app{
		exitOnError: true,
	}

	err := createCliApp(&app).Run(os.Args)

	if err != nil {
		log.Fatal("It seems there was an error: ", err)
		cli.OsExiter(1)
	}
}

type exitHandler interface {
	handle(app, error)
}

func (a *app) exitErrHandler(c *cli.Context, err error) {
	if err == nil {
		return
	}
	switch handler := err.(type) {
	case exitHandler:
		handler.handle(*a, err)
		return
	default:
		if !a.isSilent(c) {
			fmt.Println(err)
		}
		if a.exitOnError {
			cli.OsExiter(1)
		}
	}
}

func (a *app) isSilent(c *cli.Context) bool {
	if a.silent {
		return true
	}
	a.silent = c.Bool(flags.ExitCodeOnlyFlagName)
	return a.silent
}

func (a *app) logf(c *cli.Context, msg string, v ...any) {
	if a.isSilent(c) {
		return
	}
	fmt.Printf(msg+"\n", v...)
}
