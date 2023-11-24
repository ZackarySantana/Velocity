package operations

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/urfave/cli"
	"github.com/zackarysantana/velocity/src/uicli"
)

var Run = []cli.Command{
	{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("added task: ", cCtx.Args().First())

			items := []list.Item{
				uicli.SimpleItem{
					Label: "Title",
					Desc:  "desc",
				},
			}

			result, err := uicli.Run("Testing", items)

			if err != nil {
				return err
			}

			println(result)

			return nil
		},
	},
}
