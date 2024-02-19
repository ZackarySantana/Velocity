package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/cli/flags"
)

// app wraps the cli runtime to centralize common
// use-cases like error handling and logging.
type app struct {
	exitOnError bool
	silent      bool
}

// createCliApp creates the cli application with the
// given app instance.
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

// Run creates and starts the cli application.
func Run() error {
	return createCliApp(&app{}).Run(os.Args)
}

// RunAndExitOnError creates and starts the cli application
// and if an error is encountered, it will log the error and
// exit with a non-zero exit code.
func RunAndExitOnError() {
	app := app{
		exitOnError: true,
	}

	err := createCliApp(&app).Run(os.Args)

	if err != nil {
		if !app.silent {
			log.Fatal("It seems there was an error: ", err)
		}
		cli.OsExiter(1)
	}
}

// exitHandler is an interface that allows for custom
// error handling. If a cli command returns an error that
// implements this, the handle function will be called. If
// the handle function returns true, the default logging
// will not take place. If it is false, the default logging
// will take place.
type exitHandler interface {
	handle(app, error) bool
}

// exitErrHandler is the handles the exit error for the cli.
// If checks if the error implements exitHandler and calls
// the handle function if it does.
func (a *app) exitErrHandler(c *cli.Context, err error) {
	if err == nil {
		return
	}
	if handler, ok := err.(exitHandler); ok {
		if handler.handle(*a, err) {
			return
		}
	}
	if !a.isSilent(c) {
		fmt.Println(err)
	}
	if a.exitOnError {
		cli.OsExiter(1)
	}
}

// isSilent checks if the silent field has been set to true,
// and if not, it checks the cli context for the exit code
// only flag.
func (a *app) isSilent(c *cli.Context) bool {
	if a.silent {
		return true
	}
	a.silent = c.Bool(flags.ExitCodeOnlyFlagName)
	return a.silent
}

// logf is a wrapper around fmt.Printf that checks if the
// silent flag is set before logging.
func (a *app) logf(c *cli.Context, msg string, v ...any) {
	if a.isSilent(c) {
		return
	}
	fmt.Printf(msg+"\n", v...)
}
