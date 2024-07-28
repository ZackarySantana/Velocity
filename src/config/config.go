package config

import (
	"github.com/zackarysantana/velocity/src/catcher"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Tests    TestSection    `yaml:"tests"`
	Images   ImageSection   `yaml:"images"`
	Jobs     JobSection     `yaml:"jobs"`
	Routines RoutineSection `yaml:"routines"`
}

func (c *Config) Validate() error {
	catcher := catcher.New()
	catcher.Catch(c.Tests.Validate())
	catcher.Catch(c.Images.Validate())
	catcher.Catch(c.Jobs.Validate())
	catcher.Catch(c.Routines.Validate())
	return catcher.Resolve()
}

func Read(bytes []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
