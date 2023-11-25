package operations

import (
	"fmt"

	"github.com/urfave/cli"
)

var Validate = []cli.Command{
	{
		Name:      "validate",
		Aliases:   []string{"v"},
		Usage:     "validate your config file",
		ArgsUsage: "[workflow]",
		Flags: []cli.Flag{
			configFlag,
		},
		Before: CombineBefores(BeforeConfig),
		Action: func(cCtx *cli.Context) error {
			config := cCtx.String("config")
			if config == "" {
				config = "velocity.yml"
			}
			fmt.Printf("Your config '%s' is good to go!\n", config)
			return nil
		},
	},
}
