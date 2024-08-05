package flags

import (
	"context"

	"github.com/urfave/cli/v3"
)

var (
	ServerFlag = &cli.StringFlag{
		Name:  "server",
		Usage: "set the server location",
		Value: "http://localhost:8080",
	}
)

func SetServer(_ context.Context, cmd *cli.Command) error {
	cmd.Metadata["server"] = cmd.String("server")
	return nil
}

func Server(cmd *cli.Command) string {
	return cmd.Metadata["server"].(string)
}
