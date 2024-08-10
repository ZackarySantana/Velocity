package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
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
		catcher.Catch(test.error().Wrap(test.validateSyntax()))
	}
	return catcher.Resolve()
}

func (t *TestSection) validateIntegrity(c *Config) error {
	if t == nil {
		return nil
	}
	catcher := catcher.New()
	for _, test := range *t {
		catcher.Catch(test.error().Wrap(test.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (t *TestSection) error() oops.OopsErrorBuilder {
	return oops.In("test_section")
}

func (t *TestSection) ToEntities(ic service.IdCreator) []*test.Test {
	tests := make([]*test.Test, 0)
	for _, tst := range *t {
		tests = append(tests, tst.ToEntity(ic))
	}
	return tests
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

func (c *Command) error() oops.OopsErrorBuilder {
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
	catcher.ErrorWhen(t.Name == "", "name is required")
	catcher.ErrorWhen(t.Language == "" && len(t.Commands) == 0, "language or commands are required")
	catcher.ErrorWhen(t.Language != "" && len(t.Commands) > 0, "cannot specify both a language and commands")
	catcher.ErrorWhen(t.Library != "" && len(t.Commands) > 0, "cannot specify a library with commands")
	for _, cmd := range t.Commands {
		catcher.Catch(cmd.error().Wrap(cmd.validateSyntax()))
	}
	return catcher.Resolve()
}

func (t *Test) validateIntegrity(config *Config) error {
	return nil
}

func (t *Test) error() oops.OopsErrorBuilder {
	return oops.With("test_name", t.Name).With("language", t.Language).With("library", t.Library)
}

func (t *Test) ToEntity(ic service.IdCreator) *test.Test {
	cmds := make([]test.Command, len(t.Commands))
	for i, cmd := range t.Commands {
		cmds[i] = cmd.ToEntity()
	}
	return &test.Test{
		Id:        ic(),
		Name:      t.Name,
		Language:  t.Language,
		Library:   t.Library,
		Commands:  cmds,
		Directory: t.Directory,
	}
}
