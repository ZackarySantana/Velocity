package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func Run() error {
	cli := cli.App{
		Name:    "velocity",
		Version: "0.0.1",
		Usage:   "manage, run, and report on tests quickly",
		Commands: []*cli.Command{
			&Validate,
		},
		ExitErrHandler: exitErrHandler,
	}
	return cli.Run(os.Args)
}

func exitErrHandler(c *cli.Context, err error) {
	if err == nil {
		return
	}
	// TODO: telemetry?
	fmt.Println(err)
}
