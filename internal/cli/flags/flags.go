package flags

import "github.com/urfave/cli/v2"

type stringFlag struct {
	cli.StringFlag
}

func (sf stringFlag) WithDefault(d string) stringFlag {
	sf.Value = d
	return sf
}

func (sf stringFlag) Flag() *cli.StringFlag {
	return &sf.StringFlag
}

type boolFlag struct {
	cli.BoolFlag
}

func (bf boolFlag) WithDefault(d bool) boolFlag {
	bf.Value = d
	return bf
}

func (bf boolFlag) Flag() *cli.BoolFlag {
	return &bf.BoolFlag
}
