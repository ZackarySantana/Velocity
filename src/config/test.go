package config

import (
	"fmt"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type TestSection []Test

func (t *TestSection) validateSyntax() error {
	return ValidateSyntaxMany(toValidators(t))
}

func (t *TestSection) validateIntegrity(c *Config) error {
	if t == nil {
		return oops.Errorf("test section must exist and contain at least one test")
	}
	catcher := catcher.New()
	catcher.When(len(*t) == 0, "at least one test is required")
	names := make(map[string]int)
	for idx, test := range *t {
		idx2, ok := names[test.Name]
		if ok {
			catcher.Wrap(oops.Errorf("duplicate test name: %s", test.Name), "[index=%d, index_2=%d]", idx, idx2)
		}
		names[test.Name] = idx
	}
	catcher.Catch(ValidateIntegrityMany(toValidators(t), c))
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
	catcher.When(c.Shell == "" && c.Prebuilt == "", "must specify a shell or prebuilt command")
	catcher.When(c.Shell != "" && c.Prebuilt != "", "cannot specify both a shell and prebuilt command")
	catcher.When(c.Prebuilt == "" && len(c.Params) > 0, "cannot specify params without a prebuilt command")
	return catcher.Resolve()
}

func (c *Command) validateIntegrity(config *Config) error {
	return nil
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
	catcher.When(t.Name == "", "name is required")
	catcher.When(t.Language == "" && len(t.Commands) == 0, "language or commands are required")
	catcher.When(t.Language != "" && len(t.Commands) > 0, "cannot specify both a language and commands")
	catcher.When(t.Library != "" && len(t.Commands) > 0, "cannot specify a library with commands")
	for _, cmd := range t.Commands {
		catcher.Catch(cmd.validateSyntax())
	}
	return catcher.Resolve()
}

func (t *Test) validateIntegrity(config *Config) error {
	return nil
}

func (t *Test) error(i int) oops.OopsErrorBuilder {
	return oops.With(fmt.Sprintf("test_name_%d", i), t.Name).With(fmt.Sprintf("language_%d", i), t.Language).With(fmt.Sprintf("library_%d", i), t.Library).With(fmt.Sprintf("directory_%d", i), t.Directory).With(fmt.Sprintf("commands_%d", i), t.Commands)
}
