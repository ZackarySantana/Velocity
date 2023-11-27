package befores

import (
	"github.com/urfave/cli/v2"
)

func CombineBefores(beforeFuncs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, beforeFunc := range beforeFuncs {
			err := beforeFunc(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
