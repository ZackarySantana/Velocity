package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type TestSection []Test

func (t *TestSection) validateSyntax() error {
	if t == nil {
		return nil
	}
	catcher := catcher.New()
	for _, test := range *t {
		catcher.Catch(test.Error().Wrap(test.validateSyntax()))
	}
	return catcher.Resolve()
}

func (t *TestSection) validateIntegrity(c *Config) error {
	if t == nil {
		return nil
	}
	catcher := catcher.New()
	for _, test := range *t {
		catcher.Catch(test.Error().Wrap(test.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (t *TestSection) Error() oops.OopsErrorBuilder {
	return oops.Code("test_section")
}

type Command struct {
	Shell string `yaml:"shell"`

	Prebuilt string                 `yaml:"prebuilt"`
	Params   map[string]interface{} `yaml:"params"`
}

func (c *Command) validateSyntax() error {
	catcher := catcher.New()
	if c.Shell == "" && c.Prebuilt == "" {
		catcher.Error("must specify a shell or prebuilt command")
	}
	if c.Shell != "" && c.Prebuilt != "" {
		catcher.Error("cannot specify both a shell and prebuilt command")
	}
	if c.Shell != "" && len(c.Params) > 0 {
		catcher.Error("cannot specify params with a shell command")
	}
	return catcher.Resolve()
}

func (c *Command) validateIntegrity(config *Config) error {
	return nil
}

func (c *Command) Error() oops.OopsErrorBuilder {
	return oops.With("shell", c.Shell).With("prebuilt", c.Prebuilt)
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
	catcher := catcher.New()
	if t.Name == "" {
		catcher.Error("name is required")
	}
	if t.Language == "" && len(t.Commands) == 0 {
		catcher.Error("language or commands are required")
	}
	if t.Language != "" && len(t.Commands) > 0 {
		catcher.Error("cannot specify both a language and commands")
	}
	if t.Library != "" && len(t.Commands) > 0 {
		catcher.Error("cannot specify a library with commands")
	}
	for _, cmd := range t.Commands {
		catcher.Catch(cmd.Error().Wrap(cmd.validateSyntax()))
	}
	return catcher.Resolve()
}

func (t *Test) validateIntegrity(config *Config) error {
	return nil
}

func (t *Test) Error() oops.OopsErrorBuilder {
	return oops.With("test_name", t.Name).With("language", t.Language).With("library", t.Library)
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
