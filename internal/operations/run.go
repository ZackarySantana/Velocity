package operation

import (
	"fmt"

	"github.com/urfave/cli"
)

var Run = []cli.Command{
	{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("added task: ", cCtx.Args().First())
			return nil
		},
	},
}
