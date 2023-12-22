package operations

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/befores"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/internal/rc"
	"github.com/zackarysantana/velocity/src/clients/v1types"
)

var FirstTimeRegister = []*cli.Command{
	{
		Name:      "first-time-register",
		Aliases:   []string{"ftr"},
		Usage:     "register your first admin user on your server",
		ArgsUsage: "[email]",
		Flags: []cli.Flag{
			flags.Config,
		},
		Before: befores.CombineBefores(befores.Config, befores.Email, befores.VelocityClientNoAPIKey),
		Action: func(ctx *cli.Context) error {
			c, err := befores.GetConfig(ctx)
			if err != nil {
				return err
			}

			e, err := befores.GetEmail(ctx)
			if err != nil {
				return err
			}
			if e == nil {
				return cli.Exit("email is required", 1)
			}

			client, err := befores.GetVelocityClientNoAPIKey(ctx)
			if err != nil {
				return err
			}

			req := v1types.PostFirstTimeRegisterRequest{
				Email: *e,
			}
			resp, err := client.PostFirstTimeRegister(req)
			if err != nil {
				return err
			}

			rcs, err := rc.GetRuntimeConfigurations()
			if err != nil || rcs == nil {
				rcs = &rc.RuntimeConfigurations{}
			}

			r := rc.RuntimeConfiguration{
				APIKey: resp.APIKey,
				Server: *c.Config.Server,
			}
			rcs.Servers = append(rcs.Servers, r)

			err = rc.SetRuntimeConfigurations(rcs)
			if err != nil {
				return err
			}

			fmt.Printf("Success! '%s' has been registered with ID '%s'. Your api key is saved at ~/.velocityrc", *e, resp.Id)

			return nil
		},
	},
}
