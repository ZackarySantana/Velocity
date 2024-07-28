package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type TestSection []Test

func (t *TestSection) Validate() error {
	if t == nil {
		return nil
	}
	catcher := catcher.New()
	for _, test := range *t {
		catcher.Catch(validate(&test))
	}
	return catcher.Resolve()
}

type Command struct {
	Shell string `yaml:"shell"`

	Prebuilt string                 `yaml:"prebuilt"`
	Params   map[string]interface{} `yaml:"params"`
}

func (c *Command) validateSyntax() error {
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

func (c *Command) validateIntegrity(config *Config) error {
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

func (t *Test) validateSyntax() error {
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
		if err := cmd.validateSyntax(); err != nil {
			return oops.Wrapf(err, "command validation failed")
		}
	}
	return nil
}

func (t *Test) validateIntegrity(config *Config) error {
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
