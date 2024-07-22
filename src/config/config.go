package config

import "gopkg.in/yaml.v3"

type TestSection []Test

type Image struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}
type ImageSection []Image

type Job struct {
	Name   string   `yaml:"name"`
	Tests  []string `yaml:"tests"`
	Images []string `yaml:"images"`
}
type JobSection []Job

type Routine struct {
	Name string   `yaml:"name"`
	Jobs []string `yaml:"jobs"`
}
type RoutineSection []Routine

type Config struct {
	Tests    TestSection    `yaml:"tests"`
	Images   ImageSection   `yaml:"images"`
	Jobs     JobSection     `yaml:"jobs"`
	Routines RoutineSection `yaml:"routines"`
}

func Read(bytes []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
