package config

import (
	"fmt"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type TestSection []Test

func (t *TestSection) validateSyntax() error {
	if t == nil {
		return nil
	}
	catcher := catcher.New()
	for i, test := range *t {
		catcher.Catch(test.error(i).Wrap(test.validateSyntax()))
	}
	return catcher.Resolve()
}

func (t *TestSection) validateIntegrity(c *Config) error {
	if t == nil {
		return nil
	}
	catcher := catcher.New()
	for i, test := range *t {
		catcher.Catch(test.error(i).Wrap(test.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (t *TestSection) error(_ int) oops.OopsErrorBuilder {
	return oops.In("test_section")
}

type Command struct {
	Shell string `yaml:"shell"`

	Prebuilt string                 `yaml:"prebuilt"`
	Params   map[string]interface{} `yaml:"params"`
}

func (c *Command) validateSyntax() error {
	catcher := catcher.New()
	catcher.ErrorWhen(c.Shell == "" && c.Prebuilt == "", "must specify a shell or prebuilt command")
	catcher.ErrorWhen(c.Shell != "" && c.Prebuilt != "", "cannot specify both a shell and prebuilt command")
	catcher.ErrorWhen(c.Shell == "" && len(c.Params) > 0, "cannot specify params without a shell command")
	return catcher.Resolve()
}

func (c *Command) validateIntegrity(config *Config) error {
	return nil
}

func (c *Command) error(i int) oops.OopsErrorBuilder {
	return oops.With(fmt.Sprintf("shell_%d", i), c.Shell).With(fmt.Sprintf("prebuilt_%d", i), c.Prebuilt)
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
	catcher.ErrorWhen(t.Name == "", "name is required")
	catcher.ErrorWhen(t.Language == "" && len(t.Commands) == 0, "language or commands are required")
	catcher.ErrorWhen(t.Language != "" && len(t.Commands) > 0, "cannot specify both a language and commands")
	catcher.ErrorWhen(t.Library != "" && len(t.Commands) > 0, "cannot specify a library with commands")
	for i, cmd := range t.Commands {
		catcher.Catch(cmd.error(i).Wrap(cmd.validateSyntax()))
	}
	return catcher.Resolve()
}

func (t *Test) validateIntegrity(config *Config) error {
	return nil
}

func (t *Test) error(i int) oops.OopsErrorBuilder {
	return oops.With(fmt.Sprintf("test_name_%d", i), t.Name).With(fmt.Sprintf("language_%d", i), t.Language).With(fmt.Sprintf("library_%d", i), t.Library)
}
