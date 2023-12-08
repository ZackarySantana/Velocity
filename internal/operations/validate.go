package operations

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/befores"
	"github.com/zackarysantana/velocity/internal/operations/flags"
)

var Validate = []*cli.Command{
	{
		Name:      "validate",
		Aliases:   []string{"v"},
		Usage:     "validate your config file",
		ArgsUsage: "[workflow]",
		Flags: []cli.Flag{
			flags.Config,
		},
		Before: befores.CombineBefores(befores.Config),
		Action: func(ctx *cli.Context) error {
			config, err := befores.GetConfig(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("Your config '%s' is good to go!\n", config.Path)
			return nil
		},
	},
}
