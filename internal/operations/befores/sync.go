package befores

import (
	"github.com/urfave/cli/v2"
	"github.com/zackarysantana/velocity/internal/operations/flags"
)

func Sync(c *cli.Context) error {
	c.App.Metadata[flags.Sync.Name] = c.Bool(flags.Sync.Name)
	return nil
}

func GetSync(c *cli.Context) bool {
	return c.App.Metadata[flags.Sync.Name].(bool)
}
