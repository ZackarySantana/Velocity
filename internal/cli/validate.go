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
			fileName := ctx.String(flags.ConfigFlagName)
			file, err := ioutil.ReadFile(fileName)
			if err != nil {
				return fmt.Errorf("There was an error reading your file '%s': %v", fileName, err)
			}

			raw, err := config.Parse(file)
			if err != nil {
				return fmt.Errorf("There was an error parsing your file in to yaml: %v", err)
			}

			parsed, err := config.HydrateConfiguration(raw)
			if err != nil {
				return fmt.Errorf("There was an error hydrating your file: %v", err)
			}

			err = config.ValidateConfiguration(*parsed)
			if err != nil {
				return fmt.Errorf("Your file is invalid, see:\n %v", err)
			}

			fmt.Printf("Configuration (%s) is valid!\n", ctx.String("config"))

			return nil
		},
	}
	return &cmd
}
