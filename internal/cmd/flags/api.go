package flags

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/src/velocity"
)

var (
	APIFlag = &cli.StringFlag{
		Name:  "api",
		Usage: "set the api location",
		Value: "http://localhost:8080",
	}
)

func SetAPI(_ context.Context, cmd *cli.Command) error {
	cmd.Metadata[APIFlag.Name] = cmd.String(APIFlag.Name)
	return nil
}

func API(cmd *cli.Command) *velocity.APIClient {
	return velocity.NewAPIClient(cmd.Metadata[APIFlag.Name].(string))
}
