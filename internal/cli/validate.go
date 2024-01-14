package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/cli/flags"
	"github.com/zackarysantana/velocity/src/config"
)

func CreateValidate(app app) *cli.Command {
	cmd := cli.Command{
		Name:      "validate",
		Aliases:   []string{"v"},
		Usage:     "validate a configuration file",
		ArgsUsage: "[workflow]",
		Flags: []cli.Flag{
			flags.Config().Flag(),
		},
		Action: func(ctx *cli.Context) error {
			file, err := ioutil.ReadFile(ctx.String("config"))
			if err != nil {
				return err
			}

			raw, err := config.Parse(file)
			if err != nil {
				return err
			}

			parsed, err := config.HydrateConfiguration(raw)
			if err != nil {
				return err
			}

			err = config.ValidateConfiguration(*parsed)
			if err != nil {
				return err
			}

			fmt.Printf("Configuration (%s) is valid!\n", ctx.String("config"))

			return nil
		},
	}
	return &cmd
}
