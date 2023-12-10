package befores

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/flags"
	"github.com/zackarysantana/velocity/internal/rc"
	"github.com/zackarysantana/velocity/src/clients"
)

var velocityServerProviders = []func(*cli.Context) (string, error){
	velocityServerFromFlag,
	velocityServerFromConfig,
}

func VelocityClient(c *cli.Context) error {
	var velocity *string
	for _, provider := range velocityServerProviders {
		v, err := provider(c)
		if err == nil && v != "" {
			velocity = &v
		}
	}

	if velocity == nil {
		return errors.New("no velocity server was found. please include it in your yaml or as a flag option")
	}

	r, err := rc.GetRuntimeConfiguration(*velocity)
	if err != nil {
		fmt.Println("No runtime configuration found. Please provide one for ", *velocity)
		r, err = rc.AskForRuntimeConfiguration(*velocity)

		if err != nil {
			return err
		}
	}

	c.App.Metadata[flags.Velocity.Name] = *velocity
	c.App.Metadata["api_key"] = r.APIKey
	return nil
}

func velocityServerFromFlag(c *cli.Context) (string, error) {
	if c.String(flags.Velocity.Name) != "" {
		return c.String(flags.Velocity.Name), nil
	}
	return "", nil
}

func velocityServerFromConfig(c *cli.Context) (string, error) {
	config, err := GetConfig(c)
	if err != nil {
		return "", err
	}
	if config.Config.Server != nil {
		return *config.Config.Server, nil
	}
	return "", nil
}

func GetVelocityClient(c *cli.Context) (*clients.VelocityClientV1, error) {
	velocity, ok := c.App.Metadata[flags.Velocity.Name].(string)
	if !ok {
		return nil, cli.Exit("velocity not found", 1)
	}

	api_key, ok := c.App.Metadata["api_key"].(string)
	if !ok {
		return nil, cli.Exit("api key not found", 1)
	}

	client := clients.NewVelocityClientV1WithAPIKey(velocity, api_key)

	return client, nil
}
