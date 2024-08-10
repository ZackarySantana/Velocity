package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/catcher"
	"github.com/zackarysantana/velocity/src/entities"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Tests    TestSection    `yaml:"tests"`
	Images   ImageSection   `yaml:"images"`
	Jobs     JobSection     `yaml:"jobs"`
	Routines RoutineSection `yaml:"routines"`
}

func (c *Config) Validate() error {
	if c == nil {
		return oops.Errorf("config is nil")
	}
	catcher := catcher.New()
	catcher.Catch(validate(&c.Tests, c))
	catcher.Catch(validate(&c.Images, c))
	catcher.Catch(validate(&c.Jobs, c))
	catcher.Catch(validate(&c.Routines, c))
	return catcher.Resolve()
}

func (c *Config) ToEntity(ic service.IdCreator) *entities.ConfigEntity {
	ec := &entities.ConfigEntity{
		Images: c.Images.ToEntities(ic),
		Tests:  c.Tests.ToEntities(ic),
	}
	ec.Jobs = c.Jobs.ToEntities(ic, ec)
	ec.Routines = c.Routines.ToEntities(ic, ec)
	return ec
}

func Parse(bytes []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}
