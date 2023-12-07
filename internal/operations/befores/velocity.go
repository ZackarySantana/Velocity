package befores

import (
	"errors"

	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/flags"
)

var velocityServerProviders = []func(*cli.Context) (string, error){
	velocityServerFromConfig,
}

func VelocityServer(c *cli.Context) error {
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

	c.App.Metadata[flags.Velocity.Name] = *velocity
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

func GetVelocityServer(c *cli.Context) (*string, error) {
	velocity, ok := c.App.Metadata[flags.Velocity.Name].(string)
	if !ok {
		return nil, cli.Exit("velocity not found", 1)
	}
	return &velocity, nil
}
