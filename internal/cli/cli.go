package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type app struct {
	exitOnError bool
}

func createCliApp(app app) *cli.App {
	return &cli.App{
		Name:    "velocity",
		Version: "0.0.1",
		Usage:   "manage, run, and report on tests quickly",
		Commands: []*cli.Command{
			CreateValidate(app),
			CreateWorkflow(app),
		},
		ExitErrHandler: app.exitErrHandler,
	}
}

func Run() error {
	app := app{}
	return createCliApp(app).Run(os.Args)
}

func RunAndExitOnError() {
	app := app{
		exitOnError: true,
	}

	err := createCliApp(app).Run(os.Args)

	if err != nil {
		log.Fatal("It seems there was an error: ", err)
		cli.OsExiter(1)
	}
}

func (a *app) exitErrHandler(c *cli.Context, err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	if a.exitOnError {
		cli.OsExiter(1)
	}
}
