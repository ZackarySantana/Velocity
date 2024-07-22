package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type Command struct {
	Shell string `yaml:"shell"`

	Prebuilt string                 `yaml:"prebuilt"`
	Params   map[string]interface{} `yaml:"params"`
}

func (c *Command) Validate() error {
	oops := oops.With("command", *c)
	if c.Shell == "" && c.Prebuilt == "" {
		return oops.Errorf("must specify a shell or prebuilt command")
	}
	if c.Shell != "" && c.Prebuilt != "" {
		return oops.Errorf("cannot specify both a shell and prebuilt command")
	}
	if c.Shell != "" && len(c.Params) > 0 {
		return oops.Errorf("cannot specify params with a shell command")
	}
	return nil
}

func (c *Command) ToEntity() test.Command {
	return test.Command{
		Shell:    c.Shell,
		Prebuilt: c.Prebuilt,
		Params:   c.Params,
	}
}

type Test struct {
	Name string `yaml:"name"`

	Language string `yaml:"language"`
	Library  string `yaml:"library"`

	Commands []Command `yaml:"commands"`

	Directory string `yaml:"directory"`
}

func (t *Test) Validate() error {
	oops := oops.With("test", *t)
	if t.Name == "" {
		return oops.Errorf("name is required")
	}
	if t.Language == "" && len(t.Commands) == 0 {
		return oops.Errorf("language or commands are required")
	}
	if t.Language != "" && len(t.Commands) > 0 {
		return oops.Errorf("cannot specify both a language and commands")
	}
	if t.Library != "" && len(t.Commands) > 0 {
		return oops.Errorf("cannot specify a library with commands")
	}
	for _, cmd := range t.Commands {
		if err := cmd.Validate(); err != nil {
			return oops.Wrapf(err, "command validation failed")
		}
	}
	return nil
}

func (t *Test) ToEntity() *test.Test {
	cmds := make([]test.Command, len(t.Commands))
	for i, cmd := range t.Commands {
		cmds[i] = cmd.ToEntity()
	}
	return &test.Test{
		Name:      t.Name,
		Language:  t.Language,
		Library:   t.Library,
		Commands:  cmds,
		Directory: t.Directory,
	}
}
